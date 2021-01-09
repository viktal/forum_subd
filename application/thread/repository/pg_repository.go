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
			main.thread.message, main.thread.slug, main.thread.title, main.forum.slug as forum, main.thread.votes
			from main.thread
			join main.forum  on forum.forum_id = main.thread.forum_id
			join main.users u on u.user_id = forum.user_id
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
			main.thread.message, main.thread.slug, main.thread.title, main.forum.slug as forum, main.thread.votes
			from main.thread
			join main.forum  on forum.forum_id = main.thread.forum_id
			join main.users u on u.user_id = thread.user_id
			where main.thread.slug ilike '%v'`, slug)
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
		where main.thread.slug ilike ?
		returning *`, thread.Message, slug)
	} else if thread.Message == nil {
		_, err = p.db.Query(&newThread, `update main.thread
		set title = coalesce(?, title)
		where main.thread.slug ilike ?
		returning *`,  thread.Title, slug)
	} else {
		_, err = p.db.Query(&newThread, `update main.thread
		set message = coalesce(?, message),
			title = coalesce(?, title)
		where main.thread.slug ilike ?
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
	// Проверить что все родительские посты живут в том же thread
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
			select count(*) as total, bool_and(main.thread.slug ilike ?) as matched from main.post
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

	query := fmt.Sprintf(`insert into main.post values`)

	values := ""
	timeCreate := time.Now()
	for i := range posts {
		if byType == common.ID {
			id, err := strconv.Atoi(slugOrID)
			if err != nil {
				return nil, err
			}
			posts[i].ThreadID = id
			_, err = p.db.Query(&posts[i].ThreadSlug, "select main.thread.slug " +
				"from main.thread " +
				"where thread.thread_id = ?", posts[i].ThreadID)
			if err != nil {
				return nil, err
			}
			//if posts[i].ThreadSlug == "" {
			//	return nil, common.NewErr(404, "Not found")
			//}
		} else {
			posts[i].ThreadSlug = slugOrID
			_, err := p.db.Query(&posts[i].ThreadID, "select main.thread.thread_id " +
				"from main.thread " +
				"where thread.slug ilike ?", posts[i].ThreadSlug)
			if err != nil {
				return nil, err
			}
			if posts[i].ThreadID == 0 {
				return nil, common.NewErr(404, "Not found")
			}
		}


		if !posts[i].Created.IsZero() {
			posts[i].Created = timeCreate
		}

		_, err := p.db.Query(&posts[i], "select forum.forum_id, forum.slug as forum " +
			"from main.forum " +
			"join main.thread on forum.forum_id = thread.forum_id " +
			"where thread.thread_id = ?", posts[i].ThreadID)
		if err != nil {
			return nil, err
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

		values += fmt.Sprintf(`('%d', '%s', '%d', '%s', '%d', '%s', '%s', '%d', '%t', '%v')`,
			posts[i].ForumID, posts[i].Forum, posts[i].UserID, posts[i].Author, posts[i].ThreadID, posts[i].ThreadSlug,
			posts[i].Message, posts[i].Parent, posts[i].IsEdited, posts[i].Created.Format(time.RFC3339))
		if i < len(posts) - 1 {
			values += ", "
		}
	}

	query = fmt.Sprintf(`insert into main.post (forum_id, forum, user_id, author, thread_id, 
			thread, message, parent, is_edited, created) values %v returning post_id`, values)

	_, err = p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *pgStorage) GetPostsThreadByIDTree(ID int, params models.PostParams) (*[]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf(`
			WITH RECURSIVE tree
			AS
			(
			   SELECT
					post_id, forum_id, forum, thread_id, parent, created, message, author, is_edited,
					CAST (post_id AS VARCHAR (50)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
			   FROM main.post
			   WHERE post.parent = 0
			   UNION
			   SELECT
				   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
				   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(50)) , LEVEL + 1, f1.parent as sourse
			   FROM
				   tree
					   JOIN main.post f1 ON f1.parent = tree.post_id
			)
			select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created
			from tree
			where thread_id = %d
			`, ID)


	if params.Sort == common.Tree {
		if params.Desc {
			query += " order by tree.PATH desc"
		} else {
			query += " order by tree.PATH"
		}
	}

	if params.Since != nil {
		orderBy := ""
		if params.Desc {
			orderBy = " desc "
		}

		query = fmt.Sprintf(` 
				WITH RECURSIVE tree
				   AS
				   (
					   SELECT
						   post_id, forum_id, forum, thread_id, parent, created, message, author, is_edited,
						   CAST (post_id AS VARCHAR (50)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
					   FROM main.post
					   WHERE post.parent = 0
					   UNION
					   SELECT
						   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
						   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(50)) , LEVEL + 1, f1.parent as sourse
					   FROM
						   tree
					   JOIN main.post f1 ON f1.parent = tree.post_id
				   ),
				tree2 AS (
					select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created, tree.PATH
					from tree
					where thread_id = %d
					order by tree.PATH %s
				),
				tree3 AS (
					select tree2.*, ROW_NUMBER() over(order by tree2.PATH %s) AS position
					from tree2
				)
				select post_id, tree3.forum, tree3.author, tree3.thread_id, tree3.message, tree3.parent, tree3.is_edited, tree3.created
				from tree3
						 join main.thread on thread.thread_id = tree3.thread_id 
				where tree3.position > (
					select position from tree3
					where post_id = %d )
				`, ID, orderBy, orderBy, *params.Since)
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

	orderBy := ""
	if params.Desc {
		orderBy += " desc "
	}

	query := fmt.Sprintf(`
		WITH RECURSIVE base_tree AS (
			SELECT
				post_id, post.forum_id, post.forum, post.thread_id, post.parent, post.created, post.message, post.author, post.is_edited,
				CAST (post_id AS VARCHAR (500)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
			FROM main.post
			WHERE thread_id = %d and post.parent = 0
			order by post_id %s
			limit %d
		),
		 tree
		   AS
		   (
			SELECT * FROM base_tree
			   UNION
			SELECT
			   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
			   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(500)) , LEVEL + 1, tree.sourse
			FROM
			   tree
				   JOIN main.post f1 ON f1.parent = tree.post_id
		   )
		select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created
		from tree
			`, ID, orderBy, params.Limit)


	if params.Sort == common.ParentTree {
		if params.Desc {
			query += " order by tree.sourse desc, tree.PATH"
		} else {
			query += " order by tree.PATH"
		}
	}

	if params.Since != nil {
		orderBy := " order by PATH"
		if params.Desc {
			orderBy = "  order by sourse desc, PATH "
		}
		query = fmt.Sprintf(` 
				WITH RECURSIVE tree
				   AS
				   (
					   SELECT
						   post_id, forum_id, forum, thread_id, parent, created, message, author, is_edited,
						   CAST (post_id AS VARCHAR (500)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
					   FROM main.post
					   WHERE post.parent = 0
					   UNION
					   SELECT
						   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
						   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(500)) , LEVEL + 1, sourse
					   FROM
						   tree
					   JOIN main.post f1 ON f1.parent = tree.post_id
				   ),
				tree2 AS (
					select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created, tree.PATH, tree.sourse
					from tree
					where thread_id = %d
					%s
				),
				tree3 AS (
					select tree2.*, ROW_NUMBER() over(%s) AS position
					from tree2
				)
				select post_id, tree3.forum, tree3.author, tree3.thread_id, tree3.message, tree3.parent, tree3.is_edited, tree3.created
				from tree3
						 join main.thread on thread.thread_id = tree3.thread_id
				where tree3.position > (
					select position from tree3
					where post_id = %d
					)
				`, ID, orderBy, orderBy, *params.Since)
	}

	//if params.Limit != 0 {
	//	query += fmt.Sprintf(` limit %d`, params.Limit)
	//}

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

	//if params.Sort == common.Tree || params.Sort == common.ParentTree {
	//	query := fmt.Sprintf(`
	//		WITH RECURSIVE tree
	//		AS
	//		(
	//		   SELECT
	//				post_id, forum_id, forum, thread_id, parent, created, message, author, is_edited,
	//				CAST (post_id AS VARCHAR (50)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
	//		   FROM main.post
	//		   WHERE post.parent = 0
	//		   UNION
	//		   SELECT
	//			   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
	//			   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(50)) , LEVEL + 1, f1.parent as sourse
	//		   FROM
	//			   tree
	//				   JOIN main.post f1 ON f1.parent = tree.post_id
	//		)
	//		select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created
	//		from tree
	//		where thread_id = %d
	//		`, ID)
	//
	//	if params.Sort == common.Tree {
	//		if params.Desc {
	//			query += " order by tree.PATH desc"
	//		} else {
	//			query += " order by tree.PATH"
	//		}
	//	} else if params.Sort == common.ParentTree {
	//		if params.Desc {
	//			query += " order by tree.sourse desc, tree.PATH"
	//		} else {
	//			query += " order by tree.PATH"
	//		}
	//	}
	//
	//	if params.Limit != 0 {
	//		query += fmt.Sprintf(` limit %d`, params.Limit)
	//	}
	//
	//	_, err := p.db.Query(&posts, query)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return &posts, nil
	//}


	var exc struct{
		Exists bool
	}
	if posts == nil {
		_, err := p.db.Query(&exc, `
				select exists(select 1 from main.post
				where post.thread_id = ?) AS "exists"`, ID)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.Post{}, nil
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
			WITH RECURSIVE tree
			AS
			(
			   SELECT
					post_id, forum_id, forum, thread_id, parent, created, message, author, is_edited,
					CAST (post_id AS VARCHAR (50)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
			   FROM main.post
			   WHERE post.parent = 0
			   UNION
			   SELECT
				   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
				   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(50)) , LEVEL + 1, f1.parent as sourse
			   FROM
				   tree
					   JOIN main.post f1 ON f1.parent = tree.post_id
			)
			select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created
			from tree
				 join main.thread on thread.thread_id = tree.thread_id
			where thread.slug ilike '%s'
			`, slug)


	if params.Sort == common.Tree {
		if params.Desc {
			query += " order by tree.PATH desc"
		} else {
			query += " order by tree.PATH"
		}
	}

	if params.Since != nil {
		orderBy := ""
		if params.Desc {
			orderBy = " desc "
		}

		query = fmt.Sprintf(` 
				WITH RECURSIVE tree
				   AS
				   (
					   SELECT
						   post_id, forum_id, forum, thread_id, parent, created, message, author, is_edited,
						   CAST (post_id AS VARCHAR (50)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
					   FROM main.post
					   WHERE post.parent = 0
					   UNION
					   SELECT
						   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
						   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(50)) , LEVEL + 1, f1.parent as sourse
					   FROM
						   tree
					   JOIN main.post f1 ON f1.parent = tree.post_id
				   ),
				tree2 AS (
					select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created, tree.PATH
					from tree
							 join main.thread on thread.thread_id = tree.thread_id
					where thread.slug ilike '%s'
					order by tree.PATH %s
				),
				tree3 AS (
					select tree2.*, ROW_NUMBER() over(order by tree2.PATH %s) AS position
					from tree2
				)
				select post_id, tree3.forum, tree3.author, tree3.thread_id, tree3.message, tree3.parent, tree3.is_edited, tree3.created
				from tree3
						 join main.thread on thread.thread_id = tree3.thread_id 
				where tree3.position > (
					select position from tree3
					where post_id = %d )
				`, slug, orderBy, orderBy, *params.Since)
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

	orderBy := ""
	if params.Desc {
		orderBy += " desc "
	}

	query := fmt.Sprintf(`
		WITH RECURSIVE base_tree AS (
			SELECT
				post_id, post.forum_id, post.forum, post.thread_id, post.parent, post.created, post.message, post.author, post.is_edited,
				CAST (post_id AS VARCHAR (500)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
			FROM main.post
					 join main.thread on thread.thread_id = post.thread_id
			WHERE thread.slug ilike '%s' and post.parent = 0
			order by post_id %s
			limit %d
		),
		 tree
		   AS
		   (
			SELECT * FROM base_tree
			   UNION
			SELECT
			   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
			   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(500)) , LEVEL + 1, tree.sourse
			FROM
			   tree
				   JOIN main.post f1 ON f1.parent = tree.post_id
		   )
		select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created
		from tree
			`, slug, orderBy, params.Limit)


	if params.Sort == common.ParentTree {
		if params.Desc {
			query += " order by tree.sourse desc, tree.PATH"
		} else {
			query += " order by tree.PATH"
		}
	}

	if params.Since != nil {
		orderBy := " order by PATH"
		if params.Desc {
			orderBy = "  order by sourse desc, PATH "
		}
		query = fmt.Sprintf(` 
				WITH RECURSIVE tree
				   AS
				   (
					   SELECT
						   post_id, forum_id, forum, thread_id, parent, created, message, author, is_edited,
						   CAST (post_id AS VARCHAR (500)) as PATH, 0 as LEVEL, cast (post_id AS numeric) as sourse
					   FROM main.post
					   WHERE post.parent = 0
					   UNION
					   SELECT
						   f1.post_id, f1.forum_id, f1.forum, f1.thread_id, f1.parent, f1.created, f1.message, f1.author, f1.is_edited,
						   CAST ( tree.PATH ||'->'|| f1.post_id AS VARCHAR(500)) , LEVEL + 1, sourse
					   FROM
						   tree
					   JOIN main.post f1 ON f1.parent = tree.post_id
				   ),
				tree2 AS (
					select post_id, tree.forum, tree.author, tree.thread_id, tree.message, tree.parent, tree.is_edited, tree.created, tree.PATH, tree.sourse
					from tree
							 join main.thread on thread.thread_id = tree.thread_id
					where thread.slug ilike '%s'
					%s
				),
				tree3 AS (
					select tree2.*, ROW_NUMBER() over(%s) AS position
					from tree2
				)
				select post_id, tree3.forum, tree3.author, tree3.thread_id, tree3.message, tree3.parent, tree3.is_edited, tree3.created
				from tree3
						 join main.thread on thread.thread_id = tree3.thread_id
				where tree3.position > (
					select position from tree3
					where post_id = %d
					)
				`, slug, orderBy, orderBy, *params.Since)
	}

	//if params.Limit != 0 {
	//	query += fmt.Sprintf(` limit %d`, params.Limit)
	//}

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
			select post_id, post.forum, author, post.thread_id, post.message, parent, is_edited, created 
			from main.post 			
			join main.thread on thread.thread_id = post.thread_id
			where thread.slug ilike '%s'`, slug)

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
				select exists(select 1 from main.post
				join main.thread on thread.thread_id = post.thread_id
				where thread.slug ilike ?) AS "exists"`, slug)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.Post{}, nil
		}
	}

	return &posts, nil
}

func (p *pgStorage) VoteOnThreadByID(ID int, vote models.Vote) (*models.Thread, error) {
//https://stackoverflow.com/questions/6677517/update-if-different-changed
	//_, err := p.db.Query(&vote, `
	//	select * from main.vote`,
	//	vote.Nickname, ID, vote.Voice)

	_, err := p.db.Query(&vote, `
		insert into main.vote (user_id, thread_id, voice)
		values ((select user_id from main.users where nickname ilike ?), ?, '?')
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
		values ((select user_id from main.users where nickname ilike '%s'), 
				(select thread_id from main.thread where slug ilike '%s'), 
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
