package usecase

import (
	"forum/application/models"
	"forum/application/service"
	"github.com/apsdehal/go-logger"
)

type UseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  service.Repository
}

func NewUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos service.Repository) *UseCase {
	return &UseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (u UseCase) GetStatusDB() (*models.Status, error) {
	return u.repos.GetStatusDB()
}

func (u UseCase) ClearDB() error {
	return u.repos.ClearDB()
}