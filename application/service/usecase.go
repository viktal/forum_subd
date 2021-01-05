package service

import (
	"forum/application/models"
)

type UseCase interface {
	GetStatusDB() (*models.Status, error)
	ClearDB() error
}
