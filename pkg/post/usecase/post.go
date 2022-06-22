package usecase

import (
	"encoding/json"
	"forum/internal/utils/utils"
	"forum/pkg/models"
	"forum/pkg/post/repository"
	repository2 "forum/pkg/thread/repository"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type PostUsecaseInterface interface {
	ParseJsonToPosts(body io.ReadCloser) ([]models.Post, error)
	ParseJsonToPostUpdate(body io.ReadCloser) (models.PostUpdate, error)
	CreatePosts(posts models.Posts, threadId int, forumName string) ([]models.Post, int, error)
	ChangeMessage(updateMessage models.PostUpdate, id string) (models.Post, bool)
	GetParamsByQuery(query url.Values) models.FullPostParams
	GetAllInfo(params models.FullPostParams, id string) (models.FullPost, bool)
	GetPostByThread(slugOrId string, limit int, since int, sort string, desc bool) ([]models.Post, bool)
}

type PostUsecase struct {
	PostDB   repository.PostRepositoryInterface
	ThreadDB repository2.ThreadRepositoryInterface
}

func (u PostUsecase) GetPostByThread(slugOrId string, limit int, since int, sort string, desc bool) ([]models.Post, bool) {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		var ok bool
		id, ok = u.ThreadDB.GetThreadIdBySlug(slugOrId)
		if !ok {
			var re []models.Post
			return re, false
		}
	} else {
		_, ok := u.ThreadDB.GetThreadInfoById(id)
		if !ok {
			return nil, false
		}
	}

	if limit <= 0 {
		limit = 100
	}

	var ok bool
	var posts []models.Post
	switch sort {
	case "tree":
		posts, ok = u.PostDB.GetPostsTree(id, limit, since, desc)
	case "parent_tree":
		log.Println("parent_tree")
		posts, ok = u.PostDB.GetPostsParentTree(id, limit, since, desc)
	default:
		posts, ok = u.PostDB.GetAllPostByThread(id, limit, since, desc)
	}

	if !ok {
		var re []models.Post
		return re, false
	}
	if posts == nil {
		var re []models.Post
		return re, true
	}
	return posts, true
}

func (u PostUsecase) GetAllInfo(params models.FullPostParams, id string) (models.FullPost, bool) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return models.FullPost{}, false
	}

	return u.PostDB.GetAllInfo(params, intId)
}

func (u PostUsecase) GetParamsByQuery(query url.Values) models.FullPostParams {
	postParams := models.FullPostParams{
		User:   false,
		Forum:  false,
		Thread: false,
	}

	related := query.Get("related")
	if strings.Contains(related, "user") {
		postParams.User = true
	}

	if strings.Contains(related, "forum") {
		postParams.Forum = true
	}

	if strings.Contains(related, "thread") {
		postParams.Thread = true
	}

	return postParams
}

func (u PostUsecase) ChangeMessage(updateMessage models.PostUpdate, id string) (models.Post, bool) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return models.Post{}, false
	}

	return u.PostDB.ChangePost(updateMessage, intId)
}

func (u PostUsecase) ParseJsonToPostUpdate(body io.ReadCloser) (models.PostUpdate, error) {
	defer body.Close()
	var postUpdate models.PostUpdate

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&postUpdate)
	if err != nil {
		log.Println(err)
	}

	return postUpdate, err
}

func (u PostUsecase) CreatePosts(posts models.Posts, threadId int, forumName string) ([]models.Post, int, error) {
	addPosts, err := u.PostDB.AddPosts(posts, threadId, forumName)
	if err == nil {
		if addPosts == nil {
			addPosts = []models.Post{}
		}
		return addPosts, http.StatusCreated, nil
	}

	code := utils.PgxErrorCode(err)
	if code == "23503" || err.Error() == models.ErrUserUnknown {
		return []models.Post{}, http.StatusNotFound, errors.New(models.ErrUserUnknown)
	}

	return []models.Post{}, http.StatusConflict, errors.New("Parent Post is Missing")
}

func (u PostUsecase) ParseJsonToPosts(body io.ReadCloser) ([]models.Post, error) {
	defer body.Close()
	var posts []models.Post

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&posts)
	if err != nil {
		log.Println(err)
	}

	now := time.Now()
	for _, post := range posts {
		post.Created = now
	}

	return posts, err
}
