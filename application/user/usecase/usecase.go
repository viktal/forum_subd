package usecase

import (
	"forum/application/common"
	"forum/application/models"
	"forum/application/user"
	logger "github.com/apsdehal/go-logger"
)

type UseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  user.Repository
}

func NewUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos user.Repository) *UseCase {
	return &UseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (u *UseCase) GetUserProfile(id string) (*models.User, error) {
	return u.repos.GetUserByNickname(id)
}

func (u *UseCase) CreateUser(user models.User) ([]models.User, *common.Err) {
	return u.repos.CreateUser(user)
}

func (u *UseCase) UpdateUser(user models.UserUpdate) (*models.User, *common.Err) {
	us, err := u.repos.GetUserByNickname(user.Nickname)
	if err != nil {
		err := common.NewErr(404, "Not Found")
		return nil, &err
	}
	if user.Email == nil && user.Fullname == nil && user.About == nil {
		return us, nil
	}
	return u.repos.UpdateUser(user)
}
