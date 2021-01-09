package models

import "time"

type Thread struct {

	// Пользователь, создавший данную тему.
	// Required: true
	Nickname string `json:"author"`

	// Дата создания ветки на форуме.
	// Format: date-time
	CreateDate time.Time `json:"created"`

	// Форум, в котором расположена данная ветка обсуждения.
	// Read Only: true
	Forum string `json:"forum"`

	// Идентификатор ветки обсуждения.
	// Read Only: true
	ThreadID int `json:"id"`

	// Описание ветки обсуждения.
	// Required: true
	Message string `json:"message"`

	// Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL).
	// В данной структуре slug опционален и не может быть числом.
	//
	// Read Only: true
	// Pattern: ^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$
	Slug string `json:"slug"`

	// Заголовок ветки обсуждения.
	// Required: true
	Title string `json:"title"`

	// Кол-во голосов непосредственно за данное сообщение форума.
	// Read Only: true
	Votes   int32 `json:"votes"`
	ForumID int   `json:"-"`
	UserID  int   `json:"-"`
}

type ThreadUpdate struct {
	// Идентификатор ветки обсуждения.
	// Read Only: true
	ThreadID int `json:"id"`

	// Описание ветки обсуждения.
	// Required: true
	Message *string `json:"message"`


	// Заголовок ветки обсуждения.
	// Required: true
	Title *string `json:"title"`

	UserID  int   `json:"-"`
}

type ListThread []Thread

type ThreadParams struct {
	/*Идентификатор ветки обсуждения.
	  Required: true
	  In: path
	*/
	SlugOrID string
	//Флаг сортировки по убыванию.
	Desc *bool `form:"desc"`
	/*Максимальное кол-во возвращаемых записей.
	  Maximum: 10000
	  Minimum: 1
	  In: query
	  Default: 100
	*/
	Limit *int32 `form:"limit" default:"100"`
	/*Идентификатор поста, после которого будут выводиться записи
	(пост с данным идентификатором в результат не попадает).
	*/
	Since *int64 `form:"since"`
	/*Вид сортировки:

	 * flat - по дате, комментарии выводятся простым списком в порядке создания;
	 * tree - древовидный, комментарии выводятся отсортированные в дереве
	   по N штук;
	 * parent_tree - древовидные с пагинацией по родительским (parent_tree),
	   на странице N родительских комментов и все комментарии прикрепленные
	   к ним, в древвидном отображение.

	Подробности: https://park.mail.ru/blog/topic/view/1191/
	  Default: "flat"
	*/
	Sort *string `form:"sort" default:"flat"`
}
