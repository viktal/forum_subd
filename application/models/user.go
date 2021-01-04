package models

//type User struct {
//	tableName struct{} `pg:"main.users,discard_unknown_columns"`
//
//	UserID            uuid.UUID `pg:"user_id,pk,type:uuid" json:"id"`
//	UserType      string    `pg:"user_type,notnull" json:"user_type"`
//	Name          string    `pg:"name,notnull" json:"name"`
//	Surname       string    `pg:"surname" json:"surname"`
//	Email         string    `pg:"email,notnull" json:"email"`
//	PasswordHash  []byte    `pg:"password_hash,notnull" json:"-"`
//	Phone         *string   `pg:"phone" json:"phone"`
//	SocialNetwork *string   `pg:"social_network" json:"social_network"`
//}

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

	UserID int
}

type UserRequest struct {
	// Описание пользователя.
	About string `json:"about,omitempty" deepcopier:"field:About"`

	// Почтовый адрес пользователя (уникальное поле).
	// Required: true
	// Format: email
	Email string `json:"email" deepcopier:"field:Email"`

	// Полное имя пользователя.
	// Required: true
	Fullname string `json:"fullname" deepcopier:"field:Fullname"`

	// Имя пользователя (уникальное поле).
	// Данное поле допускает только латиницу, цифры и знак подчеркивания.
	// Сравнение имени регистронезависимо.
	//
	// Read Only: true
	Nickname string `json:"nickname,omitempty" deepcopier:"field:Nickname"`
}
