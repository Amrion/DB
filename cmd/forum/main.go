package main

import (
	"forum/internal/forum/middleware"
	"forum/internal/forum/repository"
	"forum/internal/utils/logger"
	delivery2 "forum/pkg/forum/delivery"
	repository2 "forum/pkg/forum/repository"
	usecase2 "forum/pkg/forum/usecase"
	"forum/pkg/models"
	delivery4 "forum/pkg/post/delivery"
	repository4 "forum/pkg/post/repository"
	usecase4 "forum/pkg/post/usecase"
	delivery5 "forum/pkg/service/delivery"
	repository5 "forum/pkg/service/repository"
	usecase5 "forum/pkg/service/usecase"
	delivery3 "forum/pkg/thread/delivery"
	repository3 "forum/pkg/thread/repository"
	usecase3 "forum/pkg/thread/usecase"
	"forum/pkg/user/delivery"
	"forum/pkg/user/repostitory"
	"forum/pkg/user/usecase"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func config() *http.Server {
	status := models.StatusInit()
	Db, err := repository.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	if Db.ProcedureRequests() != nil {
		log.Fatal(1)
		log.Fatal(err)
	}

	userUsecase := usecase.UserUsecase{DB: &repostitory.UserRepository{DB: Db.GetPostgres()}}
	forumUsecase := usecase2.ForumUsecase{DB: &repository2.ForumRepository{DB: Db.GetPostgres()}}
	threadUsecase := usecase3.ThreadUsecase{ThreadDB: &repository3.ThreadRepository{DB: Db.GetPostgres()}, ForumDB: &repository2.ForumRepository{DB: Db.GetPostgres()}}
	postUsecase := usecase4.PostUsecase{PostDB: &repository4.PostRepository{DB: Db.GetPostgres()}, ThreadDB: &repository3.ThreadRepository{DB: Db.GetPostgres()}}
	serviceUsecase := usecase5.ServiceUsecase{DB: repository5.ServiceRepository{DB: Db.GetPostgres(), Status: &status}}

	// logger
	logrus.SetFormatter(&logrus.TextFormatter{})
	mainLogger := logrus.New()
	loggerM := middleware.LoggerMiddleware{
		Logger: &logger.Logger{Logger: logrus.NewEntry(mainLogger)},
		User:   &userUsecase,
	}

	//delivery
	user := delivery.UserDeliveryStruct{Usecase: userUsecase}
	forum := delivery2.ForumDelivery{ForumUsecase: forumUsecase, ThreadUsecase: threadUsecase}
	thread := delivery3.ThreadDelivery{ThreadUsecase: threadUsecase, PostUsecase: postUsecase}
	post := delivery4.PostDelivery{Usecase: postUsecase}
	service := delivery5.ServiceDelivery{Usecase: serviceUsecase}

	//router
	mainRouter := mux.NewRouter()
	subRouter := mainRouter.PathPrefix("/api").Subrouter()
	subRouter.Use(loggerM.Middleware)

	user.SetHandlersForUsers(subRouter)
	forum.SetHandlersForForum(subRouter)
	thread.SetHandlersForThread(subRouter)
	post.SetHandlersForPost(subRouter)
	service.SetHandlersForService(subRouter)

	s := http.Server{
		Addr:    ":5000",
		Handler: mainRouter,
	}

	return &s
}

func main() {
	server := config()
	log.Println("Server Start")
	log.Fatalln(server.ListenAndServe())
}
