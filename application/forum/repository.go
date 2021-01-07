package forum

import (
	"forum/application/common"
	"forum/application/models"
)

type Repository interface {
	CreateForum(forum models.ForumCreate) (*models.ForumCreate, *common.Err)
	GetForumBySlug(slug string) (*models.Forum, error)
	GetForumByID(ID int) (*models.Forum, error)
	GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error)
	CreateThread(slugForum string, thread models.Thread) (*models.Thread, *common.Err)
	GetAllForumTreads(slugForum string, params models.ForumParams) (*[]models.Thread, error)
}
