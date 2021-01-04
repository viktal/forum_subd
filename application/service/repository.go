package service

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type Repository interface {
	GetStatusDB() (*models.Status, error)
	ClearDB() error
}
