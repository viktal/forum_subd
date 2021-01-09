package repository

import (
	"fmt"
	"forum/application/common"
	"forum/application/forum"
	"forum/application/models"
	"github.com/go-pg/pg/v9"
	"strconv"
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
	query := fmt.Sprintf(`insert into main.forum 
					(slug, title, user_id) values ('%s', '%s', '%v')`, forum.Slug, forum.Title, forum.UserID)
	_, err := p.db.Query(&forum, query)

	//var res struct {
	//	count int64
	//}
	//queryCountPost := fmt.Sprintf(`select count(*)
	//		from main.post
	//		where forum_id = '%v'
	//		group by post_id;`, forum.ForumID)
	//_, err = p.db.Query(&res, queryCountPost)
	//if err != nil {
	//	return nil, err
	//}
	//forum.Posts = res.count
	//
	//queryCountThread := fmt.Sprintf(`select count(*)
	//		from main.thread
	//		where forum_id = '%v'
	//		group by thread_id;`, forum.ForumID)
	//_, err = p.db.Query(&res, queryCountThread)
	//
	//forum.Threads = res.count


	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23505") {
			query := fmt.Sprintf(`select main.forum.title, u.nickname as user, main.forum.slug
				from main.forum
				join main.users u on u.user_id = forum.user_id
				where main.forum.slug ilike '%s'
				group by u.user_id, forum.forum_id; `, forum.Slug)
			_, err1 := p.db.Query(&forum, query)
			if err1 != nil {
				newErr := common.NewErr(500, err1.Error())
				return nil, &newErr
			} else {
				newErr := common.NewErr(409, err.Error())
				return &forum, &newErr
			}
		} else {
			newErr := common.NewErr(500, err.Error())
			return nil, &newErr
		}
	}


	return &forum, nil
}

func (p *PGRepository) GetForumBySlug(slug string) (*models.Forum, error) {
	var forum models.Forum
	forum.Slug = slug
	query := fmt.Sprintf(`select main.forum.forum_id, main.forum.slug, main.forum.title, u.nickname as user, 
				count( DISTINCT t.thread_id) threads, count( DISTINCT p.post_id) posts
				from main.forum
				join main.users u on u.user_id = forum.user_id
				left join main.post p on forum.forum_id = p.forum_id
				left join main.thread t on forum.forum_id = t.forum_id
				where main.forum.slug ilike '%s'
				group by u.user_id, forum.forum_id; `, slug)
	_, err := p.db.Query(&forum, query)
	if err != nil {
		return nil, err
	}
	return &forum, nil
}


func (p *PGRepository) GetForumByID(ID int) (*models.Forum, error) {
	var forum models.Forum
	query := fmt.Sprintf(`select main.forum.title, u.nickname as user, main.forum.slug, 
				count(t.thread_id) threads, count(p.post_id) posts
				from main.forum
				join main.users u on u.user_id = forum.user_id
				join main.post p on forum.forum_id = p.forum_id
				join main.thread t on forum.forum_id = t.forum_id
				where main.forum.forum_id = '%d'
				group by u.user_id, forum.forum_id; `, ID)
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
	//TODO params
	var req []Req
	findID := fmt.Sprintf(`select DISTINCT main.thread.user_id tu,
			array_agg(distinct main.post.user_id) as pu
			from main.forum
			join main.thread  on forum.forum_id = main.thread.forum_id
			left join main.post  on forum.forum_id = main.post.forum_id
			where main.forum.slug ilike '%s'
			group by main.thread.user_id;`, slug)

	_, err := p.db.Query(&req, findID)


	var exc struct{
		Exists bool
	}
	if req == nil {
		_, err := p.db.Query(&exc, `
				select exists(select 1 
				from main.forum
              	where forum.slug ilike ?) AS "exists"`, slug)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.User{}, nil
		} else {
			return nil, nil
		}
	}


	list := getUnique(req)
	findUsers := fmt.Sprintf(`select *
				from main.users
				where user_id in (%s)
				`, strings.Join(list, ","))

	if params.Since != nil {
		if params.Desc {
			findUsers += fmt.Sprintf(` and lower('%s')::bytea > lower(nickname)::bytea`, *params.Since)
		} else {
			findUsers += fmt.Sprintf(` and lower('%s')::bytea < lower(nickname)::bytea`, *params.Since)
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
	_, err = p.db.Query(&users, findUsers)

	if err != nil {
		return nil, err
	}

	return &users, nil
}

func getUnique(arrayStruct []Req) []string {
	keys := make(map[int]bool)
	list := []string{}
	for i, entry := range arrayStruct {
		if _, value := keys[entry.Tu]; !value {
			keys[entry.Tu] = true
			list = append(list, strconv.Itoa(entry.Tu))
		}
		for _, entry2 := range arrayStruct[i].Pu {
			if _, value := keys[entry2]; !value {
				keys[entry2] = true
				list = append(list, strconv.Itoa(entry2))
			}
		}
	}
	return list
}

func (p *PGRepository) CreateThread(slugForum string, thread models.Thread) (*models.Thread, *common.Err) {
	var UserID string
	_, err := p.db.Query(&UserID, "select user_id " +
		"from main.users where nickname = ?", thread.Nickname)
	if err != nil || UserID == "" {
		newErr := common.NewErr(404, "User not found.")
		return nil, &newErr
	}

	//var ForumID string
	_, err = p.db.Query(&thread, "select forum_id, slug as forum " +
		"from main.forum where slug ilike ?", slugForum)
	if err != nil || thread.ForumID == 0 {
		newErr := common.NewErr(404, "Forum not found")
		return nil, &newErr
	}

	if thread.CreateDate.IsZero() {
		thread.CreateDate = time.Now()
	}

	query := fmt.Sprintf(`insert into main.thread 
		(forum_id, forum, user_id, nickname, title, message, slug, create_date, votes) values 
		('%d', '%s', '%s', '%s', '%s', '%s',nullif(?, ''), ?, '%v') returning thread_id`,
		thread.ForumID, thread.Forum, UserID, thread.Nickname, thread.Title, thread.Message, thread.Votes)
	_, err = p.db.Query(&thread, query, thread.Slug, thread.CreateDate)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23505") {

			var oldThread models.Thread
			query := fmt.Sprintf(`select  main.thread.nickname, main.thread.create_date, main.thread.thread_id, 
			main.thread.message, main.thread.slug, main.thread.title, main.forum.slug as forum
			from main.thread
			join main.forum  on forum.forum_id = main.thread.forum_id
			where main.thread.slug ilike '%v'`, thread.Slug)
			_, err1 := p.db.Query(&oldThread, query)
			if err1 != nil {
				newErr := common.NewErr(500, err1.Error())
				return nil, &newErr
			} else {
				newErr := common.NewErr(409, err.Error())
				return &oldThread, &newErr
			}
		} else {
			newErr := common.NewErr(500, err.Error())
			return nil, &newErr
		}
	}
	return &thread, nil
}

func (p *PGRepository) GetAllForumTreads(slugForum string, params models.ForumParams) (*[]models.Thread, error) {
	//TODO params
	var threads []models.Thread
	query := fmt.Sprintf(`select  main.thread.*
		from main.thread
		join main.forum  on forum.forum_id = main.thread.forum_id
		where main.forum.slug ilike '%s'`, slugForum)

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
		_, err := p.db.Query(&exc, `select exists(select 1 from main.thread
				join main.forum  on forum.forum_id = main.thread.forum_id
		where main.forum.slug ilike ?) AS "exists"`, slugForum)
		if err != nil {
			return nil, err
		}
		if exc.Exists == true {
			return &[]models.Thread{}, nil
		}
	}

	return &threads, nil
}
