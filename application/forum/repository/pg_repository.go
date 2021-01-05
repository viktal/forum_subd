package repository

import (
	"fmt"
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

func (p *PGRepository) CreateForum(forum models.Forum) (*models.Forum, error) {
	query := fmt.Sprintf(`insert into main.forum 
					(slug, title, user_id) values ('%s', '%s', '%v')`, forum.Slug, forum.Title, forum.UserID)
	_, err := p.db.Query(&forum, query)
	if err != nil {
		return nil, err
	}

	var res struct {
		count int64
	}
	queryCountPost := fmt.Sprintf(`select count(*)
			from main.post
			where forum_id = '%v'
			group by post_id;`, forum.ForumID)
	_, err = p.db.Query(&res, queryCountPost)
	if err != nil {
		return nil, err
	}
	forum.Posts = res.count

	queryCountThread := fmt.Sprintf(`select count(*)
			from main.thread
			where forum_id = '%v'
			group by thread_id;`, forum.ForumID)
	_, err = p.db.Query(&res, queryCountThread)

	forum.Threads = res.count
	return &forum, nil
}

func (p *PGRepository) GetForumBySlug(slug string) (*models.Forum, error) {
	var forum models.Forum
	forum.Slug = slug
	query := fmt.Sprintf(`select main.forum.title, u.nickname as user, count(t.thread_id) threads, count(p.post_id) posts
				from main.forum
				join main.users u on u.user_id = forum.user_id
				join main.post p on forum.forum_id = p.forum_id
				join main.thread t on forum.forum_id = t.forum_id
				where main.forum.slug = '%s'
				group by u.user_id, forum.forum_id; `, slug)
	_, err := p.db.Query(&forum, query)
	if err != nil {
		return nil, err
	}
	return &forum, nil
}


func (p *PGRepository) GetForumByID(ID int) (*models.Forum, error) {
	var forum models.Forum
	query := fmt.Sprintf(`select main.forum.title, u.nickname as user, count(t.thread_id) threads, count(p.post_id) posts
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

func (p *PGRepository) GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error) {
	//TODO params
	var req []Req
	findID := fmt.Sprintf(`select DISTINCT main.thread.user_id tu,
			array_agg(distinct main.post.user_id) as pu
			from main.forum
			join main.thread  on forum.forum_id = main.thread.forum_id
			join main.post  on forum.forum_id = main.post.forum_id
			where main.forum.slug = '%s'
			group by main.thread.user_id;`, slug)


	//findID := fmt.Sprintf(`select main.post.user_id as u, main.thread.user_id as tu
	//from main.forum
	//join main.post  on forum.forum_id = main.post.forum_id
	//join main.thread  on forum.forum_id = main.thread.forum_id
    //where main.forum.slug = '%s';`, slug)
	_, err := p.db.Query(&req, findID)

	list := getUnique(req)
	findUsers := fmt.Sprintf(`select *
				from main.users
				where user_id in (%s);`, strings.Join(list, ","))

	var users []models.User
	_, err = p.db.Query(&users, findUsers)

	if err != nil {
		return nil, err
	}
	return users, nil
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

func (p *PGRepository) CreateThread(slugForum string, thread models.Thread) (*models.Thread, error) {
	var UserID string
	_, err := p.db.Query(&UserID, "select user_id " +
		"from main.users where nickname = ?", thread.Nickname)
	if err != nil {
		return nil, err
	}

	var ForumID string
	_, err = p.db.Query(&ForumID, "select forum_id " +
		"from main.forum where slug = ?", slugForum)
	if err != nil {
		return nil, err
	}
	thread.CreateDate = time.Now()

	query := fmt.Sprintf(`insert into main.thread 
		(forum_id, forum, user_id, nickname, title, message, slug, create_date, votes) values 
		('%s', '%s', '%s', '%s', '%s', '%s','%s', '%s', '%v') returning thread_id`,
		ForumID, thread.Forum, UserID, thread.Nickname, thread.Title, thread.Message,
		thread.Slug, thread.CreateDate.Format(time.RFC3339), thread.Votes)
	_, err = p.db.Query(&thread, query)
	if err != nil {
		return nil, err
	}
	thread.Forum = slugForum
	return &thread, nil
}

func (p *PGRepository) GetAllForumTreads(slugForum string, params models.ForumParams) ([]models.Thread, error) {
	//TODO params
	var threads []models.Thread
	query := fmt.Sprintf(`select  main.thread.*, 
			u.nickname as author, 
			main.forum.title as forum
		from main.thread
		join main.forum  on forum.forum_id = main.thread.forum_id
		join main.users u on u.user_id = forum.user_id
		where main.forum.slug = '%s'`, slugForum)
	_, err := p.db.Query(&threads, query)
	if err != nil {
		return nil, err
	}
	return threads, nil
}
