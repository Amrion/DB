package delivery

import (
	"forum/internal/utils/response"
	"forum/pkg/service/usecase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ServiceDeliveryInterface interface {
	CleanDB(w http.ResponseWriter, r *http.Request)
	GetFullInfo(w http.ResponseWriter, r *http.Request)
}

type ServiceDelivery struct {
	Usecase usecase.ServiceUsecaseInterface
}

func (u ServiceDelivery) SetHandlersForService(router *mux.Router) {
	router.HandleFunc("/service/clear", u.CleanDB).Methods(http.MethodPost)
	router.HandleFunc("/service/status", u.GetFullInfo).Methods(http.MethodGet)
}

func (u ServiceDelivery) CleanDB(w http.ResponseWriter, r *http.Request) {
	ok := u.Usecase.CleanDb()
	if !ok {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func (u ServiceDelivery) GetFullInfo(w http.ResponseWriter, r *http.Request) {
	status := u.Usecase.GetStatus()
	response.Process(response.LoggerFunc("GET STATUS", log.Println), response.ResponseFunc(w, http.StatusOK, status))
}
