package usecase
//
//import (
//	"fmt"
//	logger "github.com/apsdehal/go-logger"
//	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
//	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/service"
//	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
//)
//
//type UseCase struct {
//	iLog   *logger.Logger
//	errLog *logger.Logger
//	repos  service.Repository
//}
//
//func NewUseCase(iLog *logger.Logger, errLog *logger.Logger,
//	repos service.Repository) *UseCase {
//	return &UseCase{
//		iLog:   iLog,
//		errLog: errLog,
//		repos:  repos,
//	}
//}
//
//func (u *UseCase) GetThreadDetails(UserID int32, slug string) (*models.ThreadID, error) {
//	return u.repos.GetThreadDetails(UserID, slug)
//}
//
//func (u *UseCase) CreatePost(thread models.ThreadID) (*models.ThreadID, error) {
//	return u.repos.CreatePost(thread)
//}
//
//func (u *UseCase) UpdateThread(thread models.ThreadID) (*models.ThreadID, error) {
//	_, err := u.repos.GetUserByNickname(user.Nickname)
//	if err != nil {
//		err = fmt.Errorf("error get user with id : %w", err)
//		return nil, err
//	}
//
//	//var newUser models.User
//	//newUser.Nickname = user.Nickname
//	//
//	//if user.Email != "" {
//	//	newUser.Email = user.Email
//	//}
//	//if user.Fullname != "" {
//	//	newUser.Fullname = user.Fullname
//	//}
//	//if user.About != "" {
//	//	newUser.About = user.About
//	//}
//
//	newUser, err := u.repos.UpdateUser(user)
//	if err != nil {
//		err = fmt.Errorf("error in updating user with id = %s : %w", user.Nickname, err)
//		return nil, err
//	}
//
//	return newUser, nil
//	return u.repos.UpdateThread(thread)
//}
//
//func (u *UseCase) GetPostsThread(params models.ThreadParams) ([]models.ThreadID, error) {
//	return u.repos.GetPostsThread(params)
//}
//
//func (u *UseCase) VoteOnThread(vote models.Vote) (*models.ThreadID, error) {
//	return u.repos.VoteOnThread(vote)
//}
//
//
//
//func (u *UseCase) GetUserProfile(id string) (*models.User, error) {
//	userById, err := u.repos.GetUserByNickname(id)
//	if err != nil {
//		err = fmt.Errorf("error in user get by id func : %w", err)
//		return nil, err
//	}
//	return userById, nil
//}
//
//func (u *UseCase) CreateUser(user models.User) (*models.User, error) {
//	userNew, err := u.repos.CreateUser(user)
//	if err != nil {
//		if err.Error() != "user already exists" {
//			err = fmt.Errorf("error in user get by id func : %w", err)
//		}
//		return nil, err
//	}
//	return userNew, nil
//}
//
//func (u *UseCase) UpdateUser(user models.User) (*models.User, error) {
//	_, err := u.repos.GetUserByNickname(user.Nickname)
//	if err != nil {
//		err = fmt.Errorf("error get user with id : %w", err)
//		return nil, err
//	}
//
//	//var newUser models.User
//	//newUser.Nickname = user.Nickname
//	//
//	//if user.Email != "" {
//	//	newUser.Email = user.Email
//	//}
//	//if user.Fullname != "" {
//	//	newUser.Fullname = user.Fullname
//	//}
//	//if user.About != "" {
//	//	newUser.About = user.About
//	//}
//
//	newUser, err := u.repos.UpdateUser(user)
//	if err != nil {
//		err = fmt.Errorf("error in updating user with id = %s : %w", user.Nickname, err)
//		return nil, err
//	}
//
//	return newUser, nil
//}
