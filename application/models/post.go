package models

import "time"

type Post struct {
	tableName struct{} `pg:"post"`
	// Автор, написавший данное сообщение.
	// Required: true
	Author string `json:"author"`

	// Дата создания сообщения на форуме.
	// Read Only: true
	// Format: date-time
	Created time.Time `json:"created"`

	// Идентификатор форума (slug) данного сообещния.
	// Read Only: true
	Forum string `json:"forum"`

	// Идентификатор данного сообщения.
	// Read Only: true
	PostID int64 `json:"id"`

	// Истина, если данное сообщение было изменено.
	// Read Only: true
	IsEdited bool `json:"isEdited"`

	// Собственно сообщение форума.
	// Required: true
	Message string `json:"message"`

	// Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	//
	Parent int `json:"parent"`

	// Идентификатор ветви (id) обсуждения данного сообещния.
	// Read Only: true
	ThreadID int `json:"thread" pg:"thread_id"`
	Thread string `json:"-"`

	ForumID    int    `json:"-"`
	UserID     int    `json:"-"`
	ThreadSlug string `json:"-"`
}

//easyjson:json
type ListPosts []Post

type PostFull struct {
	// author
	Author *User `json:"author,omitempty"`

	// forum
	Forum *Forum `json:"forum,omitempty"`

	// post
	Post *Post `json:"post,omitempty"`

	// thread
	Thread *Thread `json:"thread,omitempty"`
}

type PostParams struct {
	Since *int `form:"since"` //Дата создания ветви обсуждения, с которой будут выводиться записи
	//Идентификатор пользователя, с которого будут выводиться пользоватли (пользователь с данным идентификатором в результат не попадает).
	Limit uint `form:"limit" default:"100"`
	Sort  string `form:"sort" default:"flat"` //Флаг сортировки.
	Desc  bool `form:"desc" default:"false"` //Флаг сортировки по убыванию.
}
