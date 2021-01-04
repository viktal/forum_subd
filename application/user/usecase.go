package user

import (
	"forum/application/models"
)

type UseCase interface {
	GetUserProfile(nickname string) (*models.UserRequest, error)
	CreateUser(user models.User) ([]models.UserRequest, error)
	UpdateUser(user models.User) (*models.UserRequest, error)
}
