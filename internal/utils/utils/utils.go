package utils

import (
	"forum/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/jackc/pgx"
	"log"
	"net/url"
)

func GetDataFromPath(param string, vars map[string]string) (string, bool) {
	data, ok := vars[param]
	return data, ok
}

func PgxErrorCode(err error) string {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return ""
	}

	return pgerr.Code
}

func ParseJsonToSearchParams(values url.Values) (models.ParamsForSearch, bool) {
	var params models.ParamsForSearch

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&params, values)

	if err != nil {
		log.Println(err)
		return params, false
	}

	if params.Limit == 0 {
		params.Limit = 100
	}

	return params, true
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func ParseJsonToGetPostsParams(values url.Values) (models.ParamsForGetPosts, bool) {
	var params models.ParamsForGetPosts

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&params, values)

	if err != nil {
		log.Println(err)
		return params, false
	}

	if params.Limit == 0 {
		params.Limit = 100
	}

	return params, true
}
