package repository

import (
	"fmt"
	"forum/application/common"
	"forum/application/forum"
	"forum/application/models"
	"github.com/go-pg/pg/v9"
	"strings"
	"time"
)

type PGRepository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) forum.Repository {
	return &PGRepository{db: db}
}

func (p *PGRepository) CreateForum(forum models.ForumCreate) (*models.ForumCreate, *common.Err) {
	query := fmt.Sprintf(`
							insert into main.forum
							(slug, title, user_id, author)
							select '%s', '%s', user_id, nickname
							from main.users
								where lower(nickname) = lower('%s')
							returning forum.slug, forum.title, forum.author, forum.forum_id`,
		forum.Slug, forum.Title, forum.Author)
	_, err := p.db.Query(&forum, query)

	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23505") {
			query := fmt.Sprintf(`
				select main.forum.title, main.forum.author as user, main.forum.slug
				from main.forum
				where lower(main.forum.slug) = lower('%s'); `, forum.Slug)
			_, err1 := p.db.Query(&forum, query)
			if err1 != nil {
				newErr := common.NewErr(409, err1.Error())
				return &forum, &newErr
			} else {
				newErr := common.NewErr(409, err.Error())
				return &forum, &newErr
			}
		} else {
			newErr := common.NewErr(404, err.Error())
			return nil, &newErr
		}
	} else if forum.ForumID == 0 {
		newErr := common.NewErr(404, "Not found")
		return nil, &newErr
	}
	return &forum, nil
}

func (p *PGRepository) GetForumBySlug(slug string) (*models.Forum, error) {
	var forum models.Forum
	forum.Slug = slug
	query := fmt.Sprintf(`select main.forum.forum_id, main.forum.slug, main.forum.title, 
				main.forum.author as user, main.forum.threads, main.forum.posts
				from main.forum
				where lower(main.forum.slug) = lower('%s')`, slug)
	_, err := p.db.Query(&forum, query)
	if err != nil {
		return nil, err
	}
	return &forum, nil
}


func (p *PGRepository) GetForumByID(ID int) (*models.Forum, error) {
	var forum models.Forum
	query := fmt.Sprintf(`select main.forum.title, main.forum.slug, 
				main.forum.author as user, main.forum.threads, main.forum.posts
				from main.forum
				where main.forum.forum_id = '%d'; `, ID)
	_, err := p.db.Query(&forum, query)
	if err != nil {
		return nil, err
	}
	return &forum, nil
}

type Req struct {
	Tu int `pq:"tu"`
	Pu []int `pg:"pu, array"`
}

func (p *PGRepository) GetAllForumUsers(slug string, params models.UserParams) (*[]models.User, error) {
	findUsers := fmt.Sprintf(`
			with thread_starters as (
				select thread_id, user_id
				from main.thread
				where lower(forum) = lower('%s')
			), all_users as (
				select post.user_id
				from thread_starters
					join main.post on post.thread_id = thread_starters.thread_id
				UNION
				SELECT user_id from thread_starters
			)
			select main.users.user_id, main.users.about, main.users.email, main.users.fullname, 
				main.users.nickname from all_users
			join main.users on all_users.user_id = main.users.user_id
			`, slug)

	if params.Since != nil {
		if params.Desc {
			findUsers += fmt.Sprintf(` where lower('%s')::bytea > nickname_byt`, *params.Since)
		} else {
			findUsers += fmt.Sprintf(` where lower('%s')::bytea < nickname_byt`, *params.Since)
		}

	}

	findUsers += ` order by lower(nickname) COLLATE "C"`

	if params.Desc {
		findUsers += " desc"
	}

	if params.Limit != 0 {
		findUsers += fmt.Sprintf(" limit %d", params.Limit)
	}

	var users []models.User
	_, err := p.db.Query(&users, findUsers)

	if err != nil {
		return nil, err
	}

	var exc struct{
		Exists bool
	}
	if users == nil {
		_, err := p.db.Query(&exc, `
				select exists(select 1
				from main.forum
	         	where lower(forum.slug) = lower(?)) AS "exists"`, slug)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.User{}, nil
		} else {
			return nil, nil
		}
	}
	return &users, nil
}

func (p *PGRepository) CreateThread(slugForum string, thread models.Thread) (*models.Thread, *common.Err) {
	if thread.CreateDate.IsZero() {
		thread.CreateDate = time.Now()
	}
	query := fmt.Sprintf(`
			insert into main.thread
			(forum, forum_id, user_id,
			nickname, title, message, slug, create_date, votes) values
			((select slug as forum from main.forum where lower(slug) = lower('%s')),
			(select forum_id as forum from main.forum where lower(slug) = lower('%s')),
			(select user_id from main.users where lower(nickname) = lower('%s')),
			'%s', '%s', '%s', nullif(?, ''), ?, %d) returning *`,
		slugForum, slugForum, thread.Nickname,
		thread.Nickname, thread.Title, thread.Message, thread.Votes)
	_, err := p.db.Query(&thread, query, thread.Slug, thread.CreateDate)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23505") {
			var oldThread models.Thread
			query := fmt.Sprintf(`select  main.thread.nickname, main.thread.create_date, main.thread.thread_id, 
			main.thread.message, main.thread.slug, main.thread.title, main.thread.forum
			from main.thread
			where lower(main.thread.slug) = lower('%v')`, thread.Slug)
			_, err1 := p.db.Query(&oldThread, query)
			if err1 != nil {
				newErr := common.NewErr(500, err1.Error())
				return nil, &newErr
			} else {
				newErr := common.NewErr(409, err.Error())
				return &oldThread, &newErr
			}
		} else {
			newErr := common.NewErr(404, err.Error())
			return nil, &newErr
		}
	}
	return &thread, nil
}

func (p *PGRepository) GetAllForumTreads(slugForum string, params models.ForumParams) (*[]models.Thread, error) {
	var threads []models.Thread
	query := fmt.Sprintf(`select main.thread.*
		from main.thread
		where lower(main.thread.forum) = lower('%s')`, slugForum)

	if !params.Since.IsZero() {
		if params.Desc {
			query += fmt.Sprintf(` and main.thread.create_date <= ? `)
		} else {
			query += fmt.Sprintf(` and main.thread.create_date >= ? `)
		}
	}

	if params.Desc {
		query += fmt.Sprintf(` order by main.thread.create_date desc`)
	} else {
		query += fmt.Sprintf(` order by main.thread.create_date`)
	}
	if params.Limit != 0 {
		query += fmt.Sprintf(` limit %d`, params.Limit)
	}

	_, err := p.db.Query(&threads, query, params.Since)
	if err != nil {
		return nil, err
	}

	var exc struct{
		Exists bool
	}
	if threads == nil {
		_, err := p.db.Query(&exc, `select exists(
				select 1 from main.thread
		where lower(main.thread.forum) = lower(?)) AS "exists"`, slugForum)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.Thread{}, nil
		}
	}

	return &threads, nil
}
