package post

import "forum/application/models"

type Repository interface {
	GetPostByID(ID int) (*models.Post, error)
	UpdatePostDetails(id int, message string) (*models.Post, error)
}
