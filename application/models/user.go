package models

type User struct {

	// Описание пользователя.
	About string `json:"about,omitempty"`

	// Почтовый адрес пользователя (уникальное поле).
	// Required: true
	// Format: email
	Email string `json:"email"`

	// Полное имя пользователя.
	// Required: true
	Fullname string `json:"fullname"`

	// Имя пользователя (уникальное поле).
	// Данное поле допускает только латиницу, цифры и знак подчеркивания.
	// Сравнение имени регистронезависимо.
	//
	// Read Only: true
	Nickname string `json:"nickname,omitempty"`

	UserID int `json:"-"`
}

//easyjson:json
type UserList []User

type UserUpdate struct {

	// Описание пользователя.
	About *string `json:"about,omitempty"`

	// Почтовый адрес пользователя (уникальное поле).
	// Required: true
	// Format: email
	Email *string `json:"email"`

	// Полное имя пользователя.
	// Required: true
	Fullname *string `json:"fullname"`

	// Имя пользователя (уникальное поле).
	// Данное поле допускает только латиницу, цифры и знак подчеркивания.
	// Сравнение имени регистронезависимо.
	//
	// Read Only: true
	Nickname string `json:"nickname,omitempty"`

	UserID int `json:"-"`
}

type UserParams struct {
	Since *string `form:"since"`
	//Идентификатор пользователя, с которого будут выводиться пользоватли (пользователь с данным идентификатором в результат не попадает).
	Limit uint `form:"limit" default:"100"`
	Desc  bool `form:"desc" default:"false"` //Флаг сортировки по убыванию.
}