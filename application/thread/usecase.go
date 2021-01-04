package thread

import (
	"forum/application/models"
)

type UseCase interface {
	GetThreadDetails(slugOrID string) (*models.Thread, error)
	UpdateThread(slugOrID string, thread models.Thread) (*models.Thread, error)
	CreatePosts(slugOrID string, posts models.ListPosts) (*models.ListPosts, error)
	GetPostsThread(params models.ThreadParams) ([]models.Thread, error)
	VoteOnThread(vote models.Vote) (*models.Thread, error)
}
