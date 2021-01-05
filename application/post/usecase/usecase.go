package usecase

import (
	"forum/application/forum"
	"forum/application/models"
	"forum/application/post"
	"forum/application/thread"
	"forum/application/user"
	"github.com/apsdehal/go-logger"
)

type UseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  post.Repository
	userRep user.Repository
	forumRep forum.Repository
	threadRep thread.Repository
}

func NewUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos post.Repository, userRep user.Repository, forumRep forum.Repository,
	threadRep thread.Repository) *UseCase {
	return &UseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
		userRep: userRep,
		forumRep: forumRep,
		threadRep: threadRep,
	}
}

func (u UseCase) GetPostDetails(ID int, related []string) (*models.PostFull, error) {
	postFull := new(models.PostFull)
	post, err := u.repos.GetPostByID(ID)
	postFull.Post = post
	if err != nil {
		return nil, err
	}

	for i := range related {
		switch related[i] {
		case "author":
			author, err := u.userRep.GetUserByID(postFull.Post.UserID)
			if err != nil {
				return nil, err
			}
			postFull.Author = author
		case "forum":
			forum, err := u.forumRep.GetForumByID(postFull.Post.ForumID)
			if err != nil {
				return nil, err
			}
			postFull.Forum = forum
		case "thread":
			thread, err := u.threadRep.GetThreadDetailsByID(postFull.Post.ThreadID)
			if err != nil {
				return nil, err
			}
			postFull.Thread = thread

		}

	}
	return postFull, err
}

func (u UseCase) UpdatePostDetails(ID int, newMessage string) (*models.Post, error) {
	return u.repos.UpdatePostDetails(ID, newMessage)
}