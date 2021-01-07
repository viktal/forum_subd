package user

import (
	"forum/application/common"
	"forum/application/models"
)

type UseCase interface {
	GetUserProfile(nickname string) (*models.User, error)
	CreateUser(user models.User) ([]models.User, *common.Err)
	UpdateUser(user models.UserUpdate) (*models.User, *common.Err)
}
