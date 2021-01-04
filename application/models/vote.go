package models

type Vote struct {

	/*Идентификатор ветки обсуждения.
	  Required: true
	  In: path
	*/
	SlugOrID string
	// Идентификатор пользователя.
	// Required: true
	Nickname string `json:"nickname"`

	// Отданный голос.
	// Required: true
	// Enum: [-1 1]
	Voice int32 `json:"voice"`
}
