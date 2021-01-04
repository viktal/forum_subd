package models

type Status struct {

	// Кол-во разделов в базе данных.
	// Required: true
	Forum int `json:"forum"`

	// Кол-во сообщений в базе данных.
	// Required: true
	Post int `json:"post"`

	// Кол-во веток обсуждения в базе данных.
	// Required: true
	Thread int `json:"thread"`

	// Кол-во пользователей в базе данных.
	// Required: true
	User int `json:"user"`
}
