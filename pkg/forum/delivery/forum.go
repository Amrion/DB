package delivery

import (
	"forum/internal/utils/response"
	"forum/internal/utils/utils"
	"forum/pkg/forum/usecase"
	"forum/pkg/models"
	usecase2 "forum/pkg/thread/usecase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ForumDeliveryInterface interface {
	CreateForum(w http.ResponseWriter, r *http.Request)
	GetForumInfo(w http.ResponseWriter, r *http.Request)
	CreateThread(w http.ResponseWriter, r *http.Request)
	GetThreadsOfForum(w http.ResponseWriter, r *http.Request)
	GetUsersOfForum(w http.ResponseWriter, r *http.Request)
}

type ForumDelivery struct {
	ForumUsecase  usecase.ForumUsecaseInterface
	ThreadUsecase usecase2.ThreadUsecaseInterface
}

func (u ForumDelivery) SetHandlersForForum(router *mux.Router) {
	router.HandleFunc("/forum/create", u.CreateForum).Methods(http.MethodPost)
	router.HandleFunc("/forum/{slug}/details", u.GetForumInfo).Methods(http.MethodGet)
	router.HandleFunc("/forum/{slug}/create", u.CreateThread).Methods(http.MethodPost)
	router.HandleFunc("/forum/{slug}/users", u.GetUsersOfForum).Methods(http.MethodGet)
	router.HandleFunc("/forum/{slug}/threads", u.GetThreadsOfForum).Methods(http.MethodGet)
}

func (d ForumDelivery) CreateForum(w http.ResponseWriter, r *http.Request) {
	forum, err := d.ForumUsecase.ParseJsonToForum(r.Body)
	if err != nil {
		return
	}

	forum, code, err := d.ForumUsecase.CreateForum(forum)
	if err != nil {
		ans := response.ErrorResponse{Err: err.Error()}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, code, ans))
		return
	}

	response.Process(response.LoggerFunc("Создан/Найден форум", log.Println), response.ResponseFunc(w, code, forum))
}

func (d ForumDelivery) GetForumInfo(w http.ResponseWriter, r *http.Request) {
	slug, ok := utils.GetDataFromPath("slug", mux.Vars(r))
	if !ok {
		return
	}

	forum, err := d.ForumUsecase.GetInfoBySlug(slug)
	if err != nil {
		ans := response.ErrorResponse{Err: err.Error()}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	response.Process(response.LoggerFunc("Вернули форум", log.Println), response.ResponseFunc(w, http.StatusOK, forum))
}

func (d ForumDelivery) CreateThread(w http.ResponseWriter, r *http.Request) {
	thread, ok := d.ThreadUsecase.GetThreadByRequest(r.Body, mux.Vars(r))
	if !ok {
		return
	}

	thread, code, err := d.ThreadUsecase.CreateThread(thread)
	if err != nil {
		ans := response.ErrorResponse{Err: err.Error()}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, code, ans))
		return
	}

	response.Process(response.LoggerFunc("Создана/Найдена ветка", log.Println), response.ResponseFunc(w, code, thread))
}

func (d ForumDelivery) GetThreadsOfForum(w http.ResponseWriter, r *http.Request) {
	slug, ok := utils.GetDataFromPath("slug", mux.Vars(r))
	if !ok {
		return
	}

	params, ok := utils.ParseJsonToSearchParams(r.URL.Query())
	if !ok {
		return
	}

	threads, ok := d.ThreadUsecase.FindThreadsByParams(slug, params)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrThreadNotfound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	response.Process(response.LoggerFunc("Return All threads By Forum", log.Println), response.ResponseFunc(w, http.StatusOK, threads))
}

func (d ForumDelivery) GetUsersOfForum(w http.ResponseWriter, r *http.Request) {
	slug, ok := utils.GetDataFromPath("slug", mux.Vars(r))
	if !ok {
		return
	}

	params, ok := utils.ParseJsonToSearchParams(r.URL.Query())
	if !ok {
		return
	}

	users, ok := d.ForumUsecase.FindUsersOfForum(slug, params)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrForumNotFound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	if users == nil {
		users = make([]models.User, 0)
	}

	response.Process(response.LoggerFunc("Return All users By Forum", log.Println), response.ResponseFunc(w, http.StatusOK, users))
}
