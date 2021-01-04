package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/thread"
	"github.com/go-pg/pg/v9"
)

func NewPgRepository(db *pg.DB) thread.Repository {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (p *pgStorage) GetThreadDetails(ID int32, slug string) (*models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) CreateThread(thread models.Thread) (*models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) UpdateThread(thread models.Thread) (*models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) GetPostsThread(params models.ThreadParams) ([]models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) VoteOnThread(vote models.Vote) (*models.Thread, error) {
	panic("implement me")
}

func (p *pgStorage) GetUserByNickname(nickname string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(`select * from main.user
							where nickname = '%s'`, nickname)

	_, err := p.db.Query(&user, query)
	if err != nil {
		return nil, err
	}
	return &user, nil

	//var newUser models.User
	//err := p.db.Model(&newUser).Where("user_id = ?", id).Select()
	//if err != nil {
	//	err = fmt.Errorf("error in select user with id: %s : error: %w", id, err)
	//	return nil, err
	//}

}

func (p *pgStorage) CreateUser(user models.User) (*models.User, error) {
	query := fmt.Sprintf(`insert into main.user 
					(nickname, email, fullname, about) values ('%s', '%s', '%s', '%s')`,
					user.Nickname, user.Email, user.Fullname, user.About)

	_, err := p.db.Query(&user, query)
	if err != nil {
		return nil, err
	}
	//_, errInsert := p.db.Model(&user).Returning("*").Insert()
	//if errInsert != nil {
	//	if isExist, err := p.db.Model(&user).Exists(); err != nil {
	//		errInsert = fmt.Errorf("error in inserting user with name: %s : error: %w", user.Nickname, err)
	//	} else if isExist {
	//		errInsert = errors.New("user already exists")
	//	}
	//	return nil, errInsert
	//}
	return &user, nil
}

func (p *pgStorage) UpdateUser(userNew models.User) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(`update main.user
				set email = '%s', fullname = '%s', about = '%s'
				where nickname = '%s'`,
				userNew.Email, userNew.Fullname, userNew.About, userNew.Nickname)

	_, err := p.db.Query(&user, query)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
