package usecase

import (
	"forum/application/models"
	"forum/application/thread"
	logger "github.com/apsdehal/go-logger"
	"strconv"
)

type UseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  thread.Repository
}

func NewUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos thread.Repository) *UseCase {
	return &UseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (u *UseCase) GetThreadDetails(slugOrID string) (*models.Thread, error) {
	var t *models.Thread
	var err error
	var id int
	if id, err = strconv.Atoi(slugOrID); err == nil {
		t, err = u.repos.GetThreadDetailsByID(id)
	} else {
		t, err = u.repos.GetThreadDetailsBySlug(slugOrID)
	}

	if err != nil {
		return nil, err
	}
	return t, nil
}

func (u *UseCase) CreateThread(threads models.ListThread) (*models.Thread, error) {
	return u.repos.CreateThread(threads)
}

func (u *UseCase) UpdateThread(slugOrID string, t models.Thread) (*models.Thread, error) {
	var newThread *models.Thread
	var err error
	var id int
	if id, err = strconv.Atoi(slugOrID); err == nil {
		newThread, err = u.repos.UpdateThreadByID(id, t)
	} else {
		newThread, err = u.repos.UpdateThreadBySlug(slugOrID, t)
	}

	if err != nil {
		return nil, err
	}
	return newThread, nil
}

func (u *UseCase) CreatePosts(slugOrID string, posts models.ListPosts) (*models.ListPosts, error) {
	var newThread *models.ListPosts
	var err error
	var id int
	if id, err = strconv.Atoi(slugOrID); err == nil {
		newThread, err = u.repos.CreatePostsByID(id, posts)
	} else {
		newThread, err = u.repos.CreatePostsBySlag(slugOrID, posts)
	}

	if err != nil {
		return nil, err
	}
	return newThread, nil
}

func (u *UseCase) GetPostsThread(params models.ThreadParams) ([]models.Thread, error) {
	return u.repos.GetPostsThread(params)
}

func (u *UseCase) VoteOnThread(vote models.Vote) (*models.Thread, error) {
	return u.repos.VoteOnThread(vote)
}
