package usecase

import (
	"encoding/json"
	"forum/internal/utils/utils"
	"forum/pkg/models"
	"forum/pkg/user/repostitory"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

type UserUsecaseInterface interface {
	ParseJsonToUser(body io.ReadCloser) (models.User, error)
	CreateUser(user models.User) ([]models.User, bool)
	GetUserByNickName(nickname string) (models.User, error)
	ChangeUser(user models.User) (models.User, int, error)
	GetUserByRequest(body io.ReadCloser, vars map[string]string) (models.User, error)
	CheckUserFields(user models.User) models.User
}

type UserUsecase struct {
	DB repostitory.UserRepositoryInterface
}

func (u UserUsecase) CheckUserFields(user models.User) models.User {
	var cleanUser models.User = models.User{Nickname: user.Nickname}
	if user.Fullname != "" {
		cleanUser.Fullname = user.Fullname
	}

	if user.About != "" {
		cleanUser.About = user.About
	}

	if user.Email != "" {
		cleanUser.Email = user.Email
	}

	return cleanUser
}

func (u UserUsecase) GetUserByRequest(body io.ReadCloser, vars map[string]string) (models.User, error) {
	user, err := u.ParseJsonToUser(body)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	nickname, ok := utils.GetDataFromPath("nickname", vars)
	if !ok {
		err = errors.New("Can't parse nickname from url")
		log.Println(err)
		return models.User{}, err
	}
	user.Nickname = nickname

	return user, nil
}

func (u UserUsecase) ChangeUser(user models.User) (models.User, int, error) {
	user, err := u.DB.ChangeUser(user)
	if err == nil {
		return user, http.StatusOK, nil
	}

	code := utils.PgxErrorCode(err)
	if code == "23503" || err == pgx.ErrNoRows {
		log.Println(err)
		return models.User{}, http.StatusNotFound, err
	}

	log.Println(err)
	return models.User{}, http.StatusConflict, err
}

func (u UserUsecase) GetUserByNickName(nickname string) (models.User, error) {
	return u.DB.GetUser(nickname)
}

func (u UserUsecase) CreateUser(user models.User) ([]models.User, bool) {
	users, ok := u.DB.AddUser(user)
	return users, ok
}

func (UserUsecase) ParseJsonToUser(body io.ReadCloser) (models.User, error) {
	defer body.Close()

	var user models.User

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&user)
	if err != nil {
	}
	return user, err
}
