package models

// Forum Информация о форуме.
type Forum struct {
	// Название форума.
	Title string `json:"title"`
	// Nickname пользователя, который отвечает за форум.
	User string `json:"user"`
	// Человекопонятный URL, уникальное поле.
	Slug string `json:"slug"`
	// Общее кол-во сообщений в данном форуме.
	Posts int64 `json:"posts"`
	// Общее кол-во ветвей обсуждения в данном форуме.
	Threads int64 `json:"threads"`
}

type ParamsForSearch struct {
	// Максимальное кол-во возвращаемых записей.
	Limit int `json:"limit"`
	// Дата создания ветви обсуждения, с которой будут выводиться записи
	Since string `json:"since"`
	// Флаг сортировки по убыванию.
	Desc bool `json:"desc"`
}

type ParamsForGetPosts struct {
	Limit int    `json:"limit"`
	Since int    `json:"since"`
	Sort  string `json:"sort"`
	Desc  bool   `json:"desc"`
}

const (
	ErrUserUnknown    = ("Такого пользователя нет!")
	ErrForumNotFound  = "Can't find forum"
	ErrPostNotFound   = "Can't find post"
	ErrThreadNotfound = "Can't find thread"
)
