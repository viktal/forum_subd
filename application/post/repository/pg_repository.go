package repository

import (
	"fmt"
	"forum/application/common"
	"forum/application/models"
	"forum/application/post"
	"github.com/go-pg/pg/v9"
)

func NewPgRepository(db *pg.DB) post.Repository {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (p pgStorage) GetPostByID(ID int) (*models.Post, error) {
	var post models.Post
	query := fmt.Sprintf(`select post_id, forum, author, user_id, forum_id, thread_id, 
								message, parent, is_edited, created 
			from main.post where post.post_id = '%d'`, ID)

	_, err := p.db.Query(&post, query)
	//fmt.Println(post.Created.Format(time.RFC3339Nano))
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p pgStorage) UpdatePostDetails(id int, message string) (*models.Post, error) {
	var post models.Post
	query := fmt.Sprintf(`update main.post
		set message = '%s',
		is_edited = main.post.is_edited or main.post.message != '%s'
		where main.post.post_id = '%v'
		returning post_id, forum, author, thread_id, message, parent, is_edited, created `, message, message, id)

	_, err := p.db.Query(&post, query)
	if err != nil {
		return nil, err
	}
	if post.PostID == 0 {
		return nil, common.NewErr(404, "Not found")
	}
	return &post, nil
}
