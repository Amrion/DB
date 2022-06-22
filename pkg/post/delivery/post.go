package delivery

import (
	"forum/internal/utils/response"
	"forum/internal/utils/utils"
	"forum/pkg/models"
	"forum/pkg/post/usecase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (u PostDelivery) SetHandlersForPost(router *mux.Router) {
	router.HandleFunc("/post/{id}/details", u.ChangePost).Methods(http.MethodPost)
	router.HandleFunc("/post/{id}/details", u.GetInfoByPost).Methods(http.MethodGet)
}

type PostDeliveryInterface interface {
	ChangePost(w http.ResponseWriter, r *http.Request)
	GetInfoByPost(w http.ResponseWriter, r *http.Request)
}

type PostDelivery struct {
	Usecase usecase.PostUsecaseInterface
}

func (u PostDelivery) GetInfoByPost(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetDataFromPath("id", mux.Vars(r))
	if !ok {
		w.WriteHeader(400)
		return
	}

	params := u.Usecase.GetParamsByQuery(r.URL.Query())

	info, ok := u.Usecase.GetAllInfo(params, id)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrPostNotFound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}

	response.Process(response.LoggerFunc("Get All info by post", log.Println), response.ResponseFunc(w, http.StatusOK, info))
}

func (u PostDelivery) ChangePost(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.GetDataFromPath("id", mux.Vars(r))
	if !ok {
		w.WriteHeader(400)
		return
	}

	updateMessage, err := u.Usecase.ParseJsonToPostUpdate(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	message, ok := u.Usecase.ChangeMessage(updateMessage, id)
	if !ok {
		ans := response.ErrorResponse{Err: models.ErrPostNotFound}
		response.Process(response.LoggerFunc(ans.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, ans))
		return
	}
	response.Process(response.LoggerFunc("Change Message", log.Println), response.ResponseFunc(w, http.StatusOK, message))
}
