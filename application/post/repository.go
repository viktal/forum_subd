package post

import "forum/application/models"

type Repository interface {
	//GetPostDetails(ID int32, related []string) (*models.PostFull, error)
	//CreateThread(thread models.Thread) (*models.Thread, error)
	//UpdateThread(thread models.Thread) (*models.Thread, error)
	//GetPostsThread(params models.ThreadParams) ([]models.Thread, error)
	//VoteOnThread(vote models.Vote) (*models.Thread, error)
	GetPostByID(ID int) (*models.Post, error)
	UpdatePostDetails(id int, message string) (*models.Post, error)
}
