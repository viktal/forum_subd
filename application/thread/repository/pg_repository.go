package repository

import (
	"fmt"
	"forum/application/common"
	"forum/application/models"
	"forum/application/thread"
	"github.com/go-pg/pg/v9"
	"strconv"
	"strings"
	"time"
)

func NewPgRepository(db *pg.DB) thread.Repository {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (p *pgStorage) GetThreadDetailsByID(ID int) (*models.Thread, error) {
	var thread models.Thread
	query := fmt.Sprintf(`select  main.thread.nickname, main.thread.create_date, main.thread.thread_id, 
			main.thread.message, main.thread.slug, main.thread.title, main.thread.forum, main.thread.votes
			from main.thread
			where main.thread.thread_id = '%v'`, ID)
	_, err := p.db.Query(&thread, query)
	if err != nil {
		return nil, err
	}
	if thread.ThreadID == 0 {
		return nil, common.NewErr(404, "Not found")
	}
	return &thread, nil
}

func (p *pgStorage) GetThreadDetailsBySlug(slug string) (*models.Thread, error) {
	var thread models.Thread
	query := fmt.Sprintf(`select  main.thread.nickname, main.thread.create_date, main.thread.thread_id, 
			main.thread.message, main.thread.slug, main.thread.title, main.thread.forum, main.thread.votes
			from main.thread
			where lower(main.thread.slug) = lower('%v')`, slug)
	_, err := p.db.Query(&thread, query)
	if err != nil {
		return nil, err
	}
	if thread.ThreadID == 0 {
		return nil, common.NewErr(404, "Not found")
	}
	return &thread, nil
}


func (p *pgStorage) UpdateThreadByID(ID int, thread models.ThreadUpdate) (*models.Thread, error) {
	var newThread models.Thread
	var err error
	if thread.Title == nil {
		_, err = p.db.Query(&newThread, `update main.thread
		set message = coalesce(?, message)
		where main.thread.thread_id = ?
		returning *`, thread.Message, ID)
	} else if thread.Message == nil {
		_, err = p.db.Query(&newThread, `update main.thread
		set title = coalesce(?, title)
		where main.thread.thread_id = ?
		returning *`, thread.Title, ID)
	} else {
		_, err = p.db.Query(&newThread, `update main.thread
		set message = coalesce(?, message),
			title = coalesce(?, title)
		where main.thread.thread_id = ?
		returning *`, thread.Message, thread.Title, ID)
	}
	if err != nil {
		return nil, err
	}
	if newThread.ThreadID == 0 {
		return nil, common.NewErr(404, "Not found")
	}
	return &newThread, nil
}

func (p *pgStorage) UpdateThreadBySlug(slug string, thread models.ThreadUpdate) (*models.Thread, error) {
	var newThread models.Thread
	var err error
	if thread.Title == nil {
		_, err = p.db.Query(&newThread, `update main.thread
		set message = coalesce(?, message)
		where lower(main.thread.slug) = lower(?)
		returning *`, thread.Message, slug)
	} else if thread.Message == nil {
		_, err = p.db.Query(&newThread, `update main.thread
		set title = coalesce(?, title)
		where lower(main.thread.slug) = lower(?)
		returning *`,  thread.Title, slug)
	} else {
		_, err = p.db.Query(&newThread, `update main.thread
		set message = coalesce(?, message),
			title = coalesce(?, title)
		where lower(main.thread.slug) = lower(?)
		returning *`, thread.Message, thread.Title, slug)
	}
	if err != nil {
		return nil, err
	}
	if newThread.ThreadID == 0 {
		return nil, common.NewErr(404, "Not found")
	}
	return &newThread, nil
}


func (p *pgStorage) CreatePosts(slugOrID string, byType string, posts models.ListPosts) (*models.ListPosts, error) {
	var parents []string
	ispresent := make(map[int]*int)
	for ind := range posts {
		parent := posts[ind].Parent
		_, ok := ispresent[parent]
		if parent != 0 && !ok {
			ispresent[parent] = nil
			parents = append(parents, strconv.Itoa(parent))
		}
	}

	var err error
	if len(parents) > 0 {
		var counts struct{
			Total int
			Matched bool
		}
		if byType == common.Slug {
			_, err = p.db.Query(&counts, `
			select count(*) as total, bool_and(lower(main.thread.slug) = lower(?)) as matched from main.post
			join main.thread on main.post.thread_id = main.thread.thread_id
			where post_id IN (?)`, slugOrID, pg.Strings(parents))
		} else {
			_, err = p.db.Query(&counts, `
			select count(*) as total, bool_and(thread_id = ?) as matched from main.post
			where post_id IN (?)`, slugOrID, pg.Strings(parents))
		}
		if err != nil {
			return nil, err
		}

		if counts.Total != len(parents) || !counts.Matched {
			err := common.NewErr(409, "Parent post was created in another thread")
			return nil, err
		}

	}

	values := ""
	timeCreate := time.Now()
	for i := range posts {
		if byType == common.ID {
			id, err := strconv.Atoi(slugOrID)
			if err != nil {
				return nil, err
			}
			posts[i].ThreadID = id
			_, err = p.db.Query(&posts[i], "select slug as thread_slug, forum_id, forum " +
				"from main.thread where thread_id = ?", posts[i].ThreadID)
			if err != nil {
				return nil, err
			}
		} else {
			posts[i].ThreadSlug = slugOrID
			_, err := p.db.Query(&posts[i], "select thread_id, forum_id, slug as thread_slug, forum " +
				"from main.thread " +
				"where lower(slug) = lower(?)", posts[i].ThreadSlug)
			if err != nil {
				return nil, err
			}
			if posts[i].ThreadID == 0 {
				return nil, common.NewErr(404, "Not found")
			}
		}


		if posts[i].Created.IsZero() {
			posts[i].Created = timeCreate
		}

		_, err = p.db.Query(&posts[i].UserID, "select main.users.user_id " +
			"from main.users " +
			"where users.nickname = ?", posts[i].Author)
		if err != nil {
			return nil, err
		}
		if posts[i].UserID == 0 {
			return nil, common.NewErr(404, "Not found")
		}

		values += fmt.Sprintf(`('%d', '%s', '%d', '%s', '%d', '%s', '%s', '%d', '%t', '%s')`,
			posts[i].ForumID, posts[i].Forum, posts[i].UserID, posts[i].Author, posts[i].ThreadID, posts[i].ThreadSlug,
			posts[i].Message, posts[i].Parent, posts[i].IsEdited, posts[i].Created.Format(time.RFC3339Nano))
		if i < len(posts) - 1 {
			values += ", "
		}
	}

	query := fmt.Sprintf(`insert into main.post (forum_id, forum, user_id, author, thread_id, 
			thread, message, parent, is_edited, created) values %s returning post_id, created`, values)

	_, err = p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *pgStorage) GetPostsThreadByIDTree(ID int, params models.PostParams) (*[]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf(`
			select post_id, forum, author, thread_id, message, parent, is_edited, created
			from main.post
			where thread_id = %d
			`, ID)

	if params.Since != nil {
		if params.Desc {
			query += fmt.Sprintf(" and path < (select path from main.post where post_id = %d) ", *params.Since)
		} else {
			query += fmt.Sprintf(" and path > (select path from main.post where post_id = %d) ", *params.Since)
		}
	}

	if params.Desc {
		query += " order by path desc, post_id desc "
	} else {
		query += " order by path, post_id "
	}

	if params.Limit != 0 {
		query += fmt.Sprintf(` limit %d`, params.Limit)
	}

	_, err := p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *pgStorage) GetPostsThreadByIDParentTree(ID int, params models.PostParams) (*[]models.Post, error) {
	var posts []models.Post

	innerQuery := fmt.Sprintf("(select post_id from main.post where thread_id = %d and parent = 0 ", ID)

	if params.Since != nil {
		if params.Desc {
			innerQuery += fmt.Sprintf(` and path[1] < (select path[1] from main.post where post_id = %d) `, *params.Since)
		} else {
			innerQuery += fmt.Sprintf(` and path[1] > (select path[1] from main.post where post_id = %d) `, *params.Since)
		}
	}

	if params.Desc {
		innerQuery += " order by post_id desc "
	} else {
		innerQuery += " order by post_id "
	}
	innerQuery += fmt.Sprintf(" limit %d) ", params.Limit)


	query := fmt.Sprintf(`select post_id, forum, author, thread_id, message, parent, is_edited, created
			  from main.post 
			  where path[1] in %s `, innerQuery)

	if params.Desc {
		query += " order by path[1] desc, path, post_id "
	} else {
		query += " order by path asc "
	}

	_, err := p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *pgStorage) GetPostsThreadByID(ID int, params models.PostParams) (*[]models.Post, error) {
	var posts []models.Post

	if params.Sort == common.Tree {
		pTree, err := p.GetPostsThreadByIDTree(ID, params)
		if err != nil {
			return nil, err
		}
		posts = *pTree
	} else if params.Sort == common.ParentTree {
		pParent, err := p.GetPostsThreadByIDParentTree(ID, params)
		if err != nil {
			return nil, err
		}
		posts = *pParent
	} else {
		query := fmt.Sprintf(`
			select post_id, forum, author, thread_id, message, parent, is_edited, created 
			from main.post where post.thread_id = %d`, ID)

		query += p.SortForGetPostsThread(params)
		_, err := p.db.Query(&posts, query)

		if err != nil {
			return nil, err
		}
	}

	var exc struct{
		Exists bool
	}
	if posts == nil {
		_, err := p.db.Query(&exc, `
				select exists(select 1 from main.thread
				where thread_id = ?) AS "exists"`, ID)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.Post{}, nil
		} else {
			return nil, nil
		}
	}

	return &posts, nil
}

func (p *pgStorage) SortForGetPostsThread(params models.PostParams) string {
	query := ""
	if params.Since != nil {
		if params.Desc {
			query += fmt.Sprintf(` and main.post.post_id < %d`, *params.Since)
		} else {
			query += fmt.Sprintf(` and main.post.post_id > %d`, *params.Since)
		}
	}
	if params.Sort == common.Flat {
		if params.Desc {
			query += fmt.Sprintf(` order by main.post.post_id desc`)
		} else {
			query += fmt.Sprintf(` order by main.post.post_id`)
		}
	}
	if params.Limit != 0 {
		query += fmt.Sprintf(` limit %d`, params.Limit)
	}
	return query
}

func (p *pgStorage) GetPostsThreadBySlugTree(slug string, params models.PostParams) (*[]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf(`
			select post_id, forum, author, thread_id, message, parent, is_edited, created
			from main.post
			where lower(thread) = lower('%s')
			`, slug)

	if params.Since != nil {
		if params.Desc {
			query += fmt.Sprintf(" and path < (select path from main.post where post_id = %d) ", *params.Since)
		} else {
			query += fmt.Sprintf(" and path > (select path from main.post where post_id = %d) ", *params.Since)
		}
	}

	if params.Desc {
		query += " order by path desc, post_id desc "
	} else {
		query += " order by path, post_id "
	}

	if params.Limit != 0 {
		query += fmt.Sprintf(` limit %d`, params.Limit)
	}

	_, err := p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *pgStorage) GetPostsThreadBySlugParentTree(slug string, params models.PostParams) (*[]models.Post, error) {
	var posts []models.Post

	innerQuery := fmt.Sprintf("(select post_id from main.post where lower(thread) = lower('%s') and parent = 0 ", slug)

	if params.Since != nil {
		if params.Desc {
			innerQuery += fmt.Sprintf(` and path[1] < (select path[1] from main.post where post_id = %d) `, *params.Since)
		} else {
			innerQuery += fmt.Sprintf(` and path[1] > (select path[1] from main.post where post_id = %d) `, *params.Since)
		}
	}

	if params.Desc {
		innerQuery += " order by post_id desc "
	} else {
		innerQuery += " order by post_id "
	}
	innerQuery += fmt.Sprintf(" limit %d) ", params.Limit)


	query := fmt.Sprintf(`select post_id, forum, author, thread_id, message, parent, is_edited, created
			  from main.post 
			  where path[1] in %s `, innerQuery)

	if params.Desc {
		query += " order by path[1] desc, path, post_id "
	} else {
		query += " order by path "
	}

	_, err := p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *pgStorage) GetPostsThreadBySlug(slug string, params models.PostParams) (*[]models.Post, error) {
	var posts []models.Post

	if params.Sort == common.Tree {
		pTree, err := p.GetPostsThreadBySlugTree(slug, params)
		if err != nil {
			return nil, err
		}
		posts = *pTree
	} else if params.Sort == common.ParentTree {
		pParent, err := p.GetPostsThreadBySlugParentTree(slug, params)
		if err != nil {
			return nil, err
		}
		posts = *pParent
	} else {
		query := fmt.Sprintf(`
			select post_id, forum, author, thread_id, message, parent, is_edited, created 
			from main.post 			
			where lower(thread) = lower('%s')`, slug)

		query += p.SortForGetPostsThread(params)
		_, err := p.db.Query(&posts, query)

		if err != nil {
			return nil, err
		}
	}

	var exc struct{
		Exists bool
	}
	if posts == nil {
		_, err := p.db.Query(&exc, `
				select exists(select 1 from main.thread
				where lower(slug) = lower(?)) AS "exists"`, slug)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.Post{}, nil
		} else {
			return nil, nil
		}
	}

	return &posts, nil
}

func (p *pgStorage) VoteOnThreadByID(ID int, vote models.Vote) (*models.Thread, error) {
	_, err := p.db.Query(&vote, `
		insert into main.vote (user_id, thread_id, voice)
		values ((select user_id from main.users where lower(nickname) = lower(?)), ?, '?')
		on conflict (user_id, thread_id) do update set voice = '?' 
		where vote.voice <> '?'`,
		vote.Nickname, ID, vote.Voice, vote.Voice, vote.Voice)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23502") {
			return nil, common.NewErr(404, "Not found")
		}
		return nil, err
	}

	return p.GetThreadDetailsByID(ID)
}

func (p *pgStorage) VoteOnThreadBySlug(slug string, vote models.Vote) (*models.Thread, error) {
	query := fmt.Sprintf(`
		insert into main.vote (user_id, thread_id, voice) 
		values ((select user_id from main.users where lower(nickname) = lower('%s')), 
				(select thread_id from main.thread where lower(slug) = lower('%s')), 
				'?')
		on conflict (user_id, thread_id) do update set voice = '?' 
		where vote.voice <> '?'`,
		vote.Nickname, slug)
	_, err := p.db.Query(&vote, query, vote.Voice, vote.Voice, vote.Voice)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23502") {
			return nil, common.NewErr(404, "Not found")
		}
		return nil, err
	}

	return p.GetThreadDetailsBySlug(slug)
}
