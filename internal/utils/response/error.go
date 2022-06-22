package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Err string `json:"message"`
}

func (e ErrorResponse) Error() string {
	ret, _ := json.Marshal(e)
	return string(ret)
}

func ResponseWithJson(w http.ResponseWriter, code int, body interface{}) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println(err)
	}
}

type (
	logfunc      func()
	responsefunc func()
)

func LoggerFunc(body interface{}, logfunc func(v ...interface{})) logfunc {
	return func() {
		logfunc(body)
	}
}

func ResponseFunc(w http.ResponseWriter, code int, body interface{}) responsefunc {
	return func() {
		ResponseWithJson(w, code, body)
	}
}

func Process(logfunc logfunc, responsefunc responsefunc) {
	logfunc()
	responsefunc()
}
