package forum

import (
	"forum/application/models"
)

type Repository interface {
	CreateForum(forum models.Forum) (*models.Forum, error)
	GetForumBySlug(slug string) (*models.Forum, error)
	GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error)
	CreateThread(slugForum string, thread models.Thread) (*models.Thread, error)
	GetAllForumTreads(slugForum string, params models.ForumParams) ([]models.Thread, error)
}
