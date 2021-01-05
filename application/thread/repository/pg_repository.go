package repository

import (
	"fmt"
	"forum/application/common"
	"forum/application/models"
	"forum/application/thread"
	"github.com/go-pg/pg/v9"
	"strconv"
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
		u.nickname, 
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


func (p *pgStorage) CreatePosts(slugOrID string, byType string, posts models.ListPosts) (*models.ListPosts, error) {
	query := fmt.Sprintf(`insert into main.post values`)

	values := ""
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
		} else {
			posts[i].ThreadSlug = slugOrID
			_, err := p.db.Query(&posts[i].ThreadID, "select main.thread.thread_id " +
				"from main.thread " +
				"where thread.slug = ?", posts[i].ThreadSlug)
			if err != nil {
				return nil, err
			}
		}

		posts[i].Created = time.Now()
		_, err := p.db.Query(&posts[i].ForumID, "select forum.forum_id " +
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

		values += fmt.Sprintf(`('%d', '%s', '%d', '%s', '%d', '%s', '%s', '%d', '%t', '%v')`,
			posts[i].ForumID, posts[i].Forum, posts[i].UserID, posts[i].Author, posts[i].ThreadID, posts[i].ThreadSlug,
			posts[i].Message, posts[i].Parent, posts[i].IsEdited, posts[i].Created.Format(time.RFC3339))
		if i < len(posts) - 1 {
			values += ", "
		}
	}

	query = fmt.Sprintf(`insert into main.post (forum_id, forum, user_id, author, thread_id, 
			thread, message, parent, is_edited, created) values %v returning post_id`, values)

	_, err := p.db.Query(&posts, query)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}


func (p *pgStorage) GetPostsThreadByID(ID int) ([]models.Post, error) {
	var posts []models.Post
	_, err := p.db.Query(&posts, `
			select post_id, forum, author, thread_id, message, parent, is_edited, created 
			from main.post where post.thread_id = ?`, ID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *pgStorage) GetPostsThreadBySlug(slug string) ([]models.Post, error) {
	var posts []models.Post
	_, err := p.db.Query(&posts, `
			select post_id, post.forum, author, post.thread_id, post.message, parent, is_edited, created 
			from main.post 			
			join main.thread on thread.thread_id = post.thread_id
			where thread.slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *pgStorage) VoteOnThreadByID(ID int, vote models.Vote) (*models.Thread, error) {
	_, err := p.db.Query(&vote, `
		insert into main.vote (user_id, thread_id, voice)
		values ((select user_id from main.users where nickname = ?), ?, '?')`,
		vote.Nickname, ID, vote.Voice)
	if err != nil {
		return nil, err
	}

	return p.GetThreadDetailsByID(ID)
}

func (p *pgStorage) VoteOnThreadBySlug(slug string, vote models.Vote) (*models.Thread, error) {

	_, err := p.db.Query(&vote, `
		insert into main.vote (user_id, thread_id, voice) 
		values ((select user_id from main.users where nickname = ?), 
				(select thread_id from main.thread where slug = ?), 
				'?')`,
		vote.Nickname, slug, vote.Voice)
	if err != nil {
		return nil, err
	}

	return p.GetThreadDetailsBySlug(slug)
}
