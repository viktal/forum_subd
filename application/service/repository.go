package service

import "forum/application/models"

type Repository interface {
	GetStatusDB() (*models.Status, error)
	ClearDB() error
}
