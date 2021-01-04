package thread

import "forum/application/models"

type Repository interface {
	GetThreadDetailsByID(ID int) (*models.Thread, error)
	GetThreadDetailsBySlug(slug string) (*models.Thread, error)

	UpdateThreadByID(ID int, thread models.Thread) (*models.Thread, error)
	UpdateThreadBySlug(slug string, thread models.Thread) (*models.Thread, error)

	CreatePostsByID(ID int, posts models.ListPosts) (*models.ListPosts, error)
	CreatePostsBySlag(Slug string, posts models.ListPosts) (*models.ListPosts, error)

	CreateThread(threads models.ListThread) (*models.Thread, error)
	GetPostsThread(params models.ThreadParams) ([]models.Thread, error)
	VoteOnThread(vote models.Vote) (*models.Thread, error)
}
