package repostitory

import (
	"forum/pkg/models"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"log"
)

type UserRepositoryInterface interface {
	AddUser(user models.User) ([]models.User, bool)
	GetUser(nickname string) (models.User, error)
	ChangeUser(user models.User) (models.User, error)
}

type UserRepository struct {
	DB *pgx.ConnPool
}

const (
	InsertUser = `INSERT INTO dbproject."User" (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`
	SelectUser = `SELECT nickname, fullname, about, email FROM dbproject."User" WHERE nickname = $1 OR email = $2`
	UpdateUser = `UPDATE dbproject."User" 
					SET fullname = COALESCE(NULLIF($1, ''), fullname), about = COALESCE(NULLIF($2, ''), about), email = COALESCE(NULLIF($3, ''), email)
					WHERE nickname = $4 
					RETURNING nickname, fullname, about, email`
	SelectUserByNick = `SELECT u.nickname, u.fullname, u.about, u.email FROM dbproject."User" u WHERE u.nickname = $1`
)

func (u *UserRepository) AddUser(user models.User) ([]models.User, bool) {
	_, err := u.DB.Exec("InsertUser", user.Nickname, user.Fullname, user.About, user.Email)
	if err == nil {
		return []models.User{user}, true
	}

	log.Println(err)
	rows, err := u.DB.Query("SelectUser", user.Nickname, user.Email)
	if err != nil {
		log.Println(err)
	}

	var users []models.User

	for rows.Next() {
		newUser := models.User{}
		err = rows.Scan(&newUser.Nickname, &newUser.Fullname, &newUser.About, &newUser.Email)
		users = append(users, newUser)
		if err != nil {
			rows.Close()
			return []models.User{}, false
		}
	}

	rows.Close()
	return users, false
}

func (u *UserRepository) ChangeUser(user models.User) (models.User, error) {
	var newUser models.User

	err := u.DB.QueryRow("UpdateUser",
		user.Fullname, user.About, user.Email, user.Nickname).
		Scan(&newUser.Nickname, &newUser.Fullname, &newUser.About, &newUser.Email)

	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	return newUser, nil
}

func (u UserRepository) GetUser(nickname string) (models.User, error) {
	var user models.User = models.User{Nickname: nickname}

	err := u.DB.QueryRow("SelectUserByNick", nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}
