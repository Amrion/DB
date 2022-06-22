package delivery

import (
	response "forum/internal/utils/response"
	"forum/internal/utils/utils"
	"forum/pkg/models"
	"forum/pkg/user/usecase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (u UserDeliveryStruct) SetHandlersForUsers(router *mux.Router) {
	router.HandleFunc("/user/{nickname}/create", u.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/{nickname}/profile", u.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/user/{nickname}/profile", u.ChangeUser).Methods(http.MethodPost)
}

type UserDeliveryInterface interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	ChangeUser(w http.ResponseWriter, r *http.Request)
}

type UserDeliveryStruct struct {
	Usecase usecase.UserUsecaseInterface
}

func (u *UserDeliveryStruct) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user, err := u.Usecase.GetUserByRequest(r.Body, mux.Vars(r))
	if err != nil {
		return
	}

	Newuser, ok := u.Usecase.CreateUser(user)
	if !ok {
		response.Process(response.LoggerFunc("Пользователь уже есть", log.Println), response.ResponseFunc(w, http.StatusConflict, Newuser))
		return
	}

	response.Process(response.LoggerFunc("Создан пользователь", log.Println), response.ResponseFunc(w, http.StatusCreated, Newuser[0]))
}

func (u UserDeliveryStruct) GetUser(w http.ResponseWriter, r *http.Request) {
	nickname, ok := utils.GetDataFromPath("nickname", mux.Vars(r))
	if !ok {
		return
	}

	user, err := u.Usecase.GetUserByNickName(nickname)
	if err != nil {
		responseError := response.ErrorResponse{Err: models.MissingUser}
		response.Process(response.LoggerFunc(responseError.Error(), log.Println), response.ResponseFunc(w, http.StatusNotFound, responseError))
		return
	}

	response.Process(response.LoggerFunc("Success get User", log.Println), response.ResponseFunc(w, http.StatusOK, user))
}

func (u UserDeliveryStruct) ChangeUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user, err := u.Usecase.GetUserByRequest(r.Body, mux.Vars(r))
	if err != nil {
		return
	}

	user = u.Usecase.CheckUserFields(user)
	user, code, err := u.Usecase.ChangeUser(user)
	if err != nil {
		responseError := response.ErrorResponse{Err: models.MissingUser}
		response.Process(response.LoggerFunc(responseError.Error(), log.Println), response.ResponseFunc(w, code, responseError))
		return
	}

	response.Process(response.LoggerFunc("Success Change User", log.Println), response.ResponseFunc(w, http.StatusOK, user))
}
