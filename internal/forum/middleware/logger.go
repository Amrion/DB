package middleware

import (
	"forum/internal/utils/logger"
	"forum/pkg/user/usecase"
	"log"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	Logger *logger.Logger
	User   *usecase.UserUsecase
}

func (l *LoggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		w.Header().Add("Content-Type", "application/json")

		log.Printf(
			"%s %s %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)

		next.ServeHTTP(w, r)
	})
}
