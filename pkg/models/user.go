package models

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}

const (
	MissingUser = "Can't find user with id #42\n"
)
