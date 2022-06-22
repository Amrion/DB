package models

import "time"

// Thread Ветка обсуждения на форуме.
type Thread struct {
	// Идентификатор ветки обсуждения.
	Id int64 `json:"id,omitempty"`
	// Заголовок ветки обсуждения.
	Title string `json:"title"`
	// Пользователь, создавший данную тему.
	Author string `json:"author"`
	// Форум, в котором расположена данная ветка обсуждения.
	Forum string `json:"forum,omitempty"`
	// Описание ветки обсуждения.
	Message string `json:"message"`
	// Кол-во голосов непосредственно за данное сообщение форума.
	Votes int64 `json:"votes"`
	// Человекопонятный URL. В данной структуре slug опционален и не может быть числом.
	Slug string `json:"slug"`
	// Дата создания ветки на форуме.
	Created time.Time `json:"created"`
}

// ThreadUpdate Сообщение для обновления ветки обсуждения на форуме. Пустые параметры остаются без изменений.
type ThreadUpdate struct {
	// Заголовок ветки обсуждения.
	Title string `json:"title"`
	// Описание ветки обсуждения.
	Message string `json:"message"`
}

// Vote Информация о голосовании пользователя.
type Vote struct {
	// Идентификатор пользователя.
	Nickname string `json:"nickname"`
	// Отданный голос.
	Voice float32 `json:"voice"`
}
