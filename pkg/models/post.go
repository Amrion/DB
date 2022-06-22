package models

import "time"

type Posts []Post

type Post struct {
	Id       int       `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type PostParams struct {
	Limit int    `json:"limit"`
	Since int    `json:"since"`
	Sort  string `json:"sort"`
	Desc  bool   `json:"desc"`
}

type FullPostParams struct {
	User   bool `json:"user"`
	Forum  bool `json:"forum"`
	Thread bool `json:"thread"`
}

type FullPost struct {
	Post   *Post   `json:"post,omitempty"`
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}
