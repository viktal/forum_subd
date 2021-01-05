package thread

import "forum/application/models"

type Repository interface {
	GetThreadDetailsByID(ID int) (*models.Thread, error)
	GetThreadDetailsBySlug(slug string) (*models.Thread, error)

	UpdateThreadByID(ID int, thread models.Thread) (*models.Thread, error)
	UpdateThreadBySlug(slug string, thread models.Thread) (*models.Thread, error)

	CreatePosts(slugOrID string, byType string, posts models.ListPosts) (*models.ListPosts, error)

	GetPostsThreadByID(ID int) ([]models.Post, error)
	GetPostsThreadBySlug(slug string) ([]models.Post, error)

	VoteOnThreadByID(ID int, vote models.Vote) (*models.Thread, error)
	VoteOnThreadBySlug(slug string, vote models.Vote) (*models.Thread, error)
}
