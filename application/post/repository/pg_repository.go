package repository

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/post"
	"github.com/go-pg/pg/v9"
)

func NewPgRepository(db *pg.DB) post.Repository {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (p pgStorage) GetPostDetails(ID int32, related []string) (*models.PostFull, error) {
	panic("implement me")
}

func (p pgStorage) CreateThread(thread models.Thread) (*models.Thread, error) {
	panic("implement me")
}

func (p pgStorage) UpdateThread(thread models.Thread) (*models.Thread, error) {
	panic("implement me")
}

func (p pgStorage) GetPostsThread(params models.ThreadParams) ([]models.Thread, error) {
	panic("implement me")
}

func (p pgStorage) VoteOnThread(vote models.Vote) (*models.Thread, error) {
	panic("implement me")
}
