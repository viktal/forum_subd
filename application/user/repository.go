package user

import "forum/application/models"

type Repository interface {
	GetUserByNickname(nickname string) (*models.User, error)
	GetUserByID(ID int) (*models.User, error)
	CreateUser(user models.User) ([]models.User, error)
	UpdateUser(user models.User) (*models.User, error)
}
