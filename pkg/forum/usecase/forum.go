package usecase

import (
	"encoding/json"
	"forum/internal/utils/utils"
	"forum/pkg/forum/repository"
	"forum/pkg/models"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

type ForumUsecaseInterface interface {
	ParseJsonToForum(body io.ReadCloser) (models.Forum, error)
	CreateForum(forum models.Forum) (models.Forum, int, error)
	GetInfoBySlug(slug string) (models.Forum, error)
	FindUsersOfForum(slug string, params models.ParamsForSearch) ([]models.User, bool)
}

type ForumUsecase struct {
	DB repository.ForumRepositoryInterface
}

func (u ForumUsecase) FindUsersOfForum(slug string, params models.ParamsForSearch) ([]models.User, bool) {
	users, ok := u.DB.FindUsers(slug, params)
	if !ok || len(users) == 0 {
		_, ok = u.DB.GetForumInfo(slug)
		if !ok {
			return nil, false
		}
		var re []models.User
		return re, true
	}

	return users, true
}

func (u ForumUsecase) GetInfoBySlug(slug string) (models.Forum, error) {
	forum, ok := u.DB.GetForumInfo(slug)
	if !ok {
		return models.Forum{}, errors.New(models.ErrForumNotFound)
	}

	return forum, nil
}

func (u ForumUsecase) CreateForum(forum models.Forum) (models.Forum, int, error) {
	log.Println(forum.User)
	var err error
	forum, err = u.DB.CreateForum(forum)
	if err == nil {
		log.Println(forum.User)
		return forum, http.StatusCreated, nil
	}

	log.Println(err)
	code := utils.PgxErrorCode(err)
	if code == "23503" { //ошибка отсутствия связей
		return models.Forum{}, http.StatusNotFound, errors.New(models.ErrUserUnknown)
	}

	forum, ok := u.DB.GetForumInfo(forum.Slug)
	if !ok {
		return models.Forum{}, http.StatusNotFound, errors.New(models.ErrUserUnknown)
	}

	return forum, http.StatusConflict, nil
}

func (u ForumUsecase) ParseJsonToForum(body io.ReadCloser) (models.Forum, error) {
	defer body.Close()
	var forum models.Forum

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&forum)

	if err != nil {
		log.Println(err)
	}

	forum.Posts = 0
	forum.Threads = 0
	return forum, err
}
