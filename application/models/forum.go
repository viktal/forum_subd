package models

import "time"

type Forum struct {

	// Общее кол-во сообщений в данном форуме.
	//
	// Read Only: true
	Posts int64 `json:"posts"`

	// Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL), уникальное поле.
	// Required: true
	// Pattern: ^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$
	Slug string `json:"slug"`

	// Общее кол-во ветвей обсуждения в данном форуме.
	//
	// Read Only: true
	Threads int64 `json:"threads"`

	// Название форума.
	// Required: true
	Title string `json:"title"`

	// Nickname пользователя, который отвечает за форум.
	// Required: true
	User    string `json:"user"`
	UserID  int    `json:"-"`
	ForumID int    `json:"-"`
}

type ForumCreate struct {

	// Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL), уникальное поле.
	// Required: true
	// Pattern: ^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$
	Slug string `json:"slug"`

	// Название форума.
	// Required: true
	Title string `json:"title"`

	// Nickname пользователя, который отвечает за форум.
	// Required: true
	User    string `json:"user"`
	UserID  int    `json:"-"`
	ForumID int    `json:"-"`
}

type ForumParams struct {
	Since time.Time `form:"since" time_format:"2006-01-02T15:04:05Z07:00"` //Дата создания ветви обсуждения, с которой будут выводиться записи
	//Идентификатор пользователя, с которого будут выводиться пользоватли (пользователь с данным идентификатором в результат не попадает).
	Limit uint `form:"limit" default:"100"`
	Desc  bool `form:"desc" default:"false"` //Флаг сортировки по убыванию.
}
