package delivery

import (
	"forum/internal/utils/response"
	"forum/internal/utils/utils"
	"forum/pkg/models"
	usecase2 "forum/pkg/post/usecase"
	"forum/pkg/thread/usecase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (u ThreadDelivery) SetHandlersForThread(router *mux.Router) {
	router.HandleFunc("/thread/{slug_or_id}/details", u.GetThreadInfo).Methods(http.MethodGet)
	router.HandleFunc("/thread/{slug_or_id}/details", u.UpdateThread).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/vote", u.VoteForThread).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/create", u.CreatePost).Methods(http.MethodPost)
	router.HandleFunc("/thread/{slug_or_id}/posts", u.GetAllPostByThread).Methods(http.MethodGet)
}

type ThreadDeliveryInterface interface {
	CreatePost(w http.ResponseWriter, r *http.Request)
	GetThreadInfo(w http.ResponseWriter, r *http.Request)
	UpdateThread(w http.ResponseWriter, r *http.Request)
	GetAllPostByThread(w http.ResponseWriter, r *http.Request)
	VoteForThread(w http.ResponseWriter, r *http.Request)
}

type ThreadDelivery struct {
	ThreadUsecase usecase.ThreadUsecaseInterface
	PostUsecase   usecase2.PostUsecaseInterface
}

func (u ThreadDelivery) GetAllPostByThread(w http.ResponseWriter, r *http.Request) {
	slugOrId, ok := utils.GetDataFromPath("slug_or_id", mux.Vars(r))
	if !ok {
		log.Println("Cant parse url")
		return
	}

	params, ok := utils.ParseJsonToGetPostsParams(r.URL.Query())
	if !ok {
		log.Println("Cant parse url params")
		return
	}

	posts, ok := u.PostUsecase.GetPostByThread(slugOrId, params.Limit, params.Since, params.Sort, params.Desc)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrThreadNotfound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	if posts == nil {
		posts = make([]models.Post, 0)
	}

	response.Process(response.LoggerFunc("Найдены посты по ветке", log.Println), response.ResponseFunc(w, http.StatusOK, posts))
}

func (u ThreadDelivery) CreatePost(w http.ResponseWriter, r *http.Request) {
	slugOrId, ok := utils.GetDataFromPath("slug_or_id", mux.Vars(r))
	if !ok {
		log.Println("Cant parse urlc")
		return
	}

	posts, err := u.PostUsecase.ParseJsonToPosts(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	thread, ok := u.ThreadUsecase.GetThreadInfo(slugOrId)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrThreadNotfound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	posts, code, err := u.PostUsecase.CreatePosts(posts, int(thread.Id), thread.Forum)
	if err != nil {
		ans := response.ErrorResponse{Err: err.Error()}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, code, ans))
		return
	}
	response.Process(response.LoggerFunc("Посты созданы", log.Println), response.ResponseFunc(w, code, posts))
}

func (u ThreadDelivery) GetThreadInfo(w http.ResponseWriter, r *http.Request) {
	slugOrId, ok := utils.GetDataFromPath("slug_or_id", mux.Vars(r))
	if !ok {
		return
	}

	thread, ok := u.ThreadUsecase.GetThreadInfo(slugOrId)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrThreadNotfound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	response.Process(response.LoggerFunc("Get info for thread", log.Println), response.ResponseFunc(w, http.StatusOK, thread))
}

func (u ThreadDelivery) VoteForThread(w http.ResponseWriter, r *http.Request) {
	slugOrId, ok := utils.GetDataFromPath("slug_or_id", mux.Vars(r))
	if !ok {
		return
	}

	vote, err := u.ThreadUsecase.ParseJsonToVote(r.Body)
	if err != nil {
		return
	}

	thread, ok := u.ThreadUsecase.SetVote(vote, slugOrId)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrThreadNotfound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	response.Process(response.LoggerFunc("Add Vote to Thread", log.Println), response.ResponseFunc(w, http.StatusOK, thread))

}

func (u ThreadDelivery) UpdateThread(w http.ResponseWriter, r *http.Request) {
	slugOrId, ok := utils.GetDataFromPath("slug_or_id", mux.Vars(r))
	if !ok {
		return
	}

	newThread, err := u.ThreadUsecase.ParseJsonToUpdateThread(r.Body)
	if err != nil {
		return
	}

	thread, ok := u.ThreadUsecase.UpdateThread(newThread, slugOrId)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrThreadNotfound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	response.Process(response.LoggerFunc("Update thread", log.Println), response.ResponseFunc(w, http.StatusOK, thread))

}
