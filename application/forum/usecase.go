package forum

import (
	"forum/application/common"
	"forum/application/models"
)

type UseCase interface {
	CreateForum(template models.ForumCreate) (*models.ForumCreate, *common.Err)
	CreateThread(slugForum string, thread models.Thread) (*models.Thread, *common.Err)
	GetForumBySlug(slug string) (*models.Forum, error)
	GetAllForumTreads(slugForum string, params models.ForumParams) (*[]models.Thread, error)
	GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error)
}
