package usecase

import (
	"forum/application/common"
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


func (u *UseCase) UpdateThread(slugOrID string, t models.ThreadUpdate) (*models.Thread, error) {
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
	var _ int
	if _, err = strconv.Atoi(slugOrID); err == nil {
		newThread, err = u.repos.CreatePosts(slugOrID, common.ID, posts)
	} else {
		newThread, err = u.repos.CreatePosts(slugOrID, common.Slug, posts)
	}

	if err != nil {
		return nil, err
	}
	return newThread, nil
}

func (u *UseCase) GetPostsThread(slugOrID string, params models.PostParams) (*[]models.Post, error) {
	var posts *[]models.Post
	var err error
	var id int
	if id, err = strconv.Atoi(slugOrID); err == nil {
		posts, err = u.repos.GetPostsThreadByID(id, params)
	} else {
		posts, err = u.repos.GetPostsThreadBySlug(slugOrID, params)
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (u *UseCase) VoteOnThread(slugOrID string, vote models.Vote) (*models.Thread, error) {
	var thread *models.Thread
	var err error
	var id int
	if id, err = strconv.Atoi(slugOrID); err == nil {
		thread, err = u.repos.VoteOnThreadByID(id, vote)
	} else {
		thread, err = u.repos.VoteOnThreadBySlug(slugOrID, vote)
	}
	if err != nil {
		return nil, err
	}
	return thread, nil
}
