package usecase

import (
	"fmt"
	logger "github.com/apsdehal/go-logger"
	"forum/application/models"
	"forum/application/user"
	"github.com/ulule/deepcopier"
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

func (u *UseCase) GetUserProfile(id string) (*models.UserRequest, error) {
	userById, err := u.repos.GetUserByNickname(id)
	if err != nil {
		return nil, err
	}

	req := &models.UserRequest{}
	err = deepcopier.Copy(userById).To(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *UseCase) CreateUser(user models.User) ([]models.UserRequest, error) {
	userNew, err := u.repos.CreateUser(user)
	if err != nil {
		return nil, err
	}

	var listReq []models.UserRequest
	for i := range userNew {
		req := &models.UserRequest{}
		err = deepcopier.Copy(userNew[i]).To(req)
		if err != nil {
			return nil, err
		}
		listReq = append(listReq, *req)
	}
	return listReq, nil
}

func (u *UseCase) UpdateUser(user models.User) (*models.UserRequest, error) {
	_, err := u.repos.GetUserByNickname(user.Nickname)
	if err != nil {
		err = fmt.Errorf("%w, code 404", err)
		return nil, err
	}
	newUser, err := u.repos.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	req := &models.UserRequest{}
	err = deepcopier.Copy(newUser).To(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
