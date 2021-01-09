package repository

import (
	"forum/application/models"
	"forum/application/service"
	"github.com/go-pg/pg/v9"
)

func NewPgRepository(db *pg.DB) service.Repository {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (p pgStorage) GetStatusDB() (*models.Status, error) {
	var status models.Status
	_, err := p.db.Query(&status, "select count(main.users.user_id) as user from main.users")
	if err != nil {
		return nil, err
	}

	_, err = p.db.Query(&status, "" +
		"select count( * ) as forum, " +
		"sum( main.forum.threads) as thread, " +
		"sum( main.forum.posts) as post " +
		"from main.forum")
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (p pgStorage) ClearDB() error {
	_, err := p.db.Query(nil, "" +
		"delete from main.vote;" +
		"delete from main.post;" +
		"delete from main.thread;" +
		"delete from main.forum;" +
		"delete from main.users;")
	return err
}
