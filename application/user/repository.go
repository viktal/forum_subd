package user

import (
	"forum/application/common"
	"forum/application/models"
)

type Repository interface {
	GetUserByNickname(nickname string) (*models.User, error)
	GetUserByID(ID int) (*models.User, error)
	CreateUser(user models.User) ([]models.User, *common.Err)
	UpdateUser(user models.UserUpdate) (*models.User, *common.Err)
}
