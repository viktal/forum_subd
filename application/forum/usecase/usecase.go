package usecase

import (
	"forum/application/common"
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

func (u *UseCase) CreateForum(forum models.ForumCreate) (*models.ForumCreate, *common.Err) {
	user, err := u.strgUser.GetUserByNickname(forum.User)
	if err != nil {
		newErr := common.NewErr(404, err.Error())
		return nil, &newErr
	}
	forum.UserID = user.UserID
	forum.User = user.Nickname

	return u.strg.CreateForum(forum)
}

func (u *UseCase) CreateThread(slugForum string, thread models.Thread) (*models.Thread, *common.Err) {
	return u.strg.CreateThread(slugForum, thread)
}

func (u *UseCase) GetForumBySlug(slug string) (*models.Forum, error) {
	return u.strg.GetForumBySlug(slug)
}

func (u *UseCase) GetAllForumTreads(slug string, params models.ForumParams) (*[]models.Thread, error) {
	return u.strg.GetAllForumTreads(slug, params)
}

func (u *UseCase) GetAllForumUsers(slug string, params models.UserParams) (*[]models.User, error) {
	return u.strg.GetAllForumUsers(slug, params)
}
