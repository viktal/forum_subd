package post

import (
	"forum/application/models"
)

type UseCase interface {
	GetPostDetails(ID int, related []string) (*models.PostFull, error)
	UpdatePostDetails(ID int, newMessage string) (*models.Post, error)
}
