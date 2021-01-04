package forum

import (
	"forum/application/models"
)

type UseCase interface {
	CreateForum(template models.Forum) (*models.Forum, error)
	CreateThread(slugForum string, thread models.Thread) (*models.Thread, error)
	GetForumBySlug(slug string) (*models.Forum, error)
	GetAllForumTreads(slugForum string, params models.ForumParams) ([]models.Thread, error)
	GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error)
}
