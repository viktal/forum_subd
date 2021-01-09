package thread

import (
	"forum/application/models"
)

type UseCase interface {
	GetThreadDetails(slugOrID string) (*models.Thread, error)
	UpdateThread(slugOrID string, thread models.ThreadUpdate) (*models.Thread, error)
	CreatePosts(slugOrID string, posts models.ListPosts) (*models.ListPosts, error)
	GetPostsThread(slugOrID string, params models.PostParams) (*[]models.Post, error)
	VoteOnThread(slugOrID string, vote models.Vote) (*models.Thread, error)
}
