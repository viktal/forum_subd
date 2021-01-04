package usecase

import (
	"forum/application/forum"
	"forum/application/models"
	"forum/application/user"
	"github.com/apsdehal/go-logger"
)

type UseCase struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        forum.Repository
	strgUser    user.Repository
}

func NewUseCase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	strg forum.Repository,
	strgUser user.Repository) forum.UseCase {
	usecase := UseCase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		strg:        strg,
		strgUser:    strgUser,
	}
	return &usecase
}

func (u *UseCase) CreateForum(forum models.Forum) (*models.Forum, error) {
	user, err := u.strgUser.GetUserByNickname(forum.User)
	if err != nil {
		return nil, err
	}
	forum.UserID = user.UserID

	return u.strg.CreateForum(forum)
}

func (u *UseCase) CreateThread(slugForum string, thread models.Thread) (*models.Thread, error) {
	return u.strg.CreateThread(slugForum, thread)
}

func (u *UseCase) GetForumBySlug(slug string) (*models.Forum, error) {
	return u.strg.GetForumBySlug(slug)
}

func (u *UseCase) GetAllForumTreads(slug string, params models.ForumParams) ([]models.Thread, error) {
	return u.strg.GetAllForumTreads(slug, params)
}

func (u *UseCase) GetAllForumUsers(slug string, params models.ForumParams) ([]models.User, error) {
	return u.strg.GetAllForumUsers(slug, params)
}
