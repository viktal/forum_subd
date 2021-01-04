package post

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type UseCase interface {
	GetPostDetails(ID int, related []string) (*models.PostFull, error)
	UpdatePostDetails(ID int, newMessage string) (*models.Post, error)
}
