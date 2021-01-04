package repository

import (
	"fmt"
	"forum/application/models"
	"forum/application/thread"
	"github.com/go-pg/pg/v9"
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
	query := fmt.Sprintf(`select  main.thread.*, 
		u.nickname as author, 
		main.forum.title as forum
			from main.thread
			join main.forum  on forum.forum_id = main.thread.forum_id
			join main.users u on u.user_id = forum.user_id
			where main.thread.thread_id = '%v'`, ID)
	_, err := p.db.Query(&thread, query)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (p *pgStorage) GetThreadDetailsBySlug(slug string) (*models.Thread, error) {
	var thread models.Thread
	query := fmt.Sprintf(`select  main.thread.*, 
		u.nickname as author, 
		main.forum.title as forum
			from main.thread
			join main.forum  on forum.forum_id = main.thread.forum_id
			join main.users u on u.user_id = forum.user_id
			where main.thread.slug = '%v'`, slug)
	_, err := p.db.Query(&thread, query)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}


func (p *pgStorage) UpdateThreadByID(ID int, thread models.Thread) (*models.Thread, error) {
	//TODO Returning *
	query := fmt.Sprintf(`update main.thread
		set message = coalesce('%s', message),
			title = coalesce('%s', title)
		where main.thread.thread_id = '%v'
		returning *`, thread.Message, thread.Title, ID)

	_, err := p.db.Query(&thread, query)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (p *pgStorage) UpdateThreadBySlug(slug string, thread models.Thread) (*models.Thread, error) {
	//TODO Returning *
	query := fmt.Sprintf(`update main.thread
		set message = coalesce('%s', message),
			title = coalesce('%s', title)
		where main.thread.slug = '%s'
		returning *`, thread.Message, thread.Title, slug)

	_, err := p.db.Query(&thread, query)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}


func (p *pgStorage) CreatePostsByID(ID int, posts models.ListPosts) (*models.ListPosts, error) {
	query := fmt.Sprintf(`insert into main.post values`)

	values := ""
	for _, post := range posts {
		post.ThreadID = ID
		post.Created = time.Now()
		if post.Forum != "" {
			_, err := p.db.Query(&post.ForumID, "select main.forum.forum_id " +
				"from main.forum " +
				"where forum.slug = ?", post.Forum)
			if err != nil {
				return nil, err
			}
		}

		_, err := p.db.Query(&post.UserID, "select main.users.user_id " +
			"from main.users " +
			"where users.nickname = ?", post.Author)
		if err != nil {
			return nil, err
		}

		_, err = p.db.Query(&post.ThreadSlug, "select main.thread.slug " +
			"from main.thread " +
			"where thread.thread_id = ?", post.ThreadID)
		if err != nil {
			return nil, err
		}

		values += fmt.Sprintf(`('%d', '%s', '%d', '%s', '%d', '%s', '%s', '%d', '%t', '%v')`,
			post.ForumID, post.Forum, post.UserID, post.Author, post.ThreadID, post.ThreadSlug, post.Message,
			post.Parent, post.IsEdited, post.Created.Format(time.RFC3339))
	}

	query = fmt.Sprintf(`insert into main.post (forum_id, forum, user_id, author, thread_id, 
			thread, message, parent, is_edited, created) values %v`, values)

	_, err := p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *pgStorage) CreatePostsBySlag(Slug string, posts models.ListPosts) (*models.ListPosts, error) {
	panic("implement me")
}



func (p *pgStorage) CreateThread(threads models.ListThread) (*models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) UpdateThread(thread models.Thread) (*models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) GetPostsThread(params models.ThreadParams) ([]models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) VoteOnThread(vote models.Vote) (*models.Thread, error) {
	panic("implement me")
}
