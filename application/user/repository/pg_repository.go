package repository

import (
	"fmt"
	"forum/application/common"
	"forum/application/models"
	"forum/application/user"
	"github.com/go-pg/pg/v9"
	"strings"
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
							where nickname ilike '%s'`, nickname)

	_, err := p.db.Query(&user, query)
	if err != nil || user.Nickname == ""{
		err = fmt.Errorf("%w, code 404", err)
		return nil, err
	}
	return &user, nil
}

func (p *pgStorage) CreateUser(user models.User) ([]models.User, *common.Err) {
	var listUsers []models.User
	//INSERT INTO ....
	//VALUES ......
	//ON CONFLICT DO NOTHING RETURNING email, nickname;

	query := fmt.Sprintf(`insert into main.users
					(nickname, email, fullname, about)
					values ('%s', '%s', '%s', '%s')`,
					user.Nickname, user.Email, user.Fullname, user.About)

	//TODO: Пользователь уже присутсвует в базе данных. Возвращает данные ранее созданных пользователей с тем же nickname-ом иои email-ом.

	_, err := p.db.Query(&user, query)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23505") {
			query := fmt.Sprintf(`
				select * from main.users 
				where nickname ilike '%s' or email ilike '%s'`, user.Nickname, user.Email)
			_, err1 := p.db.Query(&listUsers, query)
			if err1 != nil {
				newErr := common.NewErr(500, err1.Error())
				return nil, &newErr
			} else {
				newErr := common.NewErr(409, err.Error())
				return listUsers, &newErr
			}
		} else {
			newErr := common.NewErr(500, err.Error())
			return nil, &newErr
		}
	}
	listUsers = append(listUsers, user)
	return listUsers, nil
}

func (p *pgStorage) UpdateUser(userNew models.UserUpdate) (*models.User, *common.Err) {
	var user models.User
	_, err := p.db.Query(&user, `update main.users
				set 
				email = COALESCE(?, email),
				fullname = COALESCE(?, fullname), 
				about = COALESCE(?, about)
				where nickname = ?
				returning *`,
		userNew.Email, userNew.Fullname, userNew.About, userNew.Nickname)
	if user.UserID == 0 {
		newErr := common.NewErr(409, "Already exist")
		return nil, &newErr
	}

	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR #23505"){
			newErr := common.NewErr(409, err.Error())
			return nil, &newErr
		}
		newErr := common.NewErr(404, err.Error())
		return nil, &newErr
	}
	return &user, nil
}
