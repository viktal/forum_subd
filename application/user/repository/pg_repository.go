package repository

import (
	"fmt"
	"forum/application/models"
	"forum/application/user"
	"github.com/go-pg/pg/v9"
)

func NewPgRepository(db *pg.DB) user.Repository {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (p *pgStorage) GetUserByID(ID int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(`select * from main.users
							where user_id = '%d'`, ID)

	_, err := p.db.Query(&user, query)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *pgStorage) GetUserByNickname(nickname string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(`select * from main.users
							where nickname = '%s'`, nickname)

	_, err := p.db.Query(&user, query)
	if err != nil || user.Nickname == ""{
		err = fmt.Errorf("%w, code 404", err)
		return nil, err
	}
	return &user, nil
}

func (p *pgStorage) CreateUser(user models.User) ([]models.User, error) {
	var listUsers []models.User
	//INSERT INTO ....
	//VALUES ......
	//ON CONFLICT DO NOTHING RETURNING email, nickname;

	query := fmt.Sprintf(`insert into main.users
					(nickname, email, fullname, about)
					values ('%s', '%s', '%s', '%s') 
				ON CONFLICT DO NOTHING RETURNING nickname, email, fullname, about;`,
					user.Nickname, user.Email, user.Fullname, user.About)

	//TODO: Пользователь уже присутсвует в базе данных. Возвращает данные ранее созданных пользователей с тем же nickname-ом иои email-ом.

	oldUser, err := p.db.Query(&user, query)
	if err != nil {
		fmt.Print(oldUser)
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
	listUsers = append(listUsers, user)
	return listUsers, nil
}

func (p *pgStorage) UpdateUser(userNew models.User) (*models.User, error) {
	query := fmt.Sprintf(`update main.users
				set email = '%s', fullname = '%s', about = '%s'
				where nickname = '%s'`,
				userNew.Email, userNew.Fullname, userNew.About, userNew.Nickname)

	_, err := p.db.Query(&userNew, query)
	if err != nil {
		return nil, err
	}
	return &userNew, nil
}
