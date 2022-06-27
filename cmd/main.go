package main

import (
	"context"
	"fmt"
	"log"
	"tp-db-project/internal/delivery"
	"tp-db-project/repository"
	"tp-db-project/repository/postgres"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"
)

func InitDb() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), "postgres://root:root@localhost:5432/dbms")
	if err != nil {
		log.Fatalln("conn error:", err)
	}

	return pool
}

func InitRepos(db *pgxpool.Pool) (repository.UserRepo, repository.ForumRepo, repository.ThreadRepo, repository.PostRepo, repository.VoteRepo, repository.StatusRepo) {
	userRepo := postgres.InitUserRepository(db)
	forumRepo := postgres.InitForumRepository(db)
	threadRepo := postgres.InitThreadRepository(db)
	postRepo := postgres.InitPostRepository(db)
	voteRepo := postgres.InitVoteRepository(db)
	statusRepo := postgres.InitStatusRepository(db)

	return userRepo, forumRepo, threadRepo, postRepo, voteRepo, statusRepo
}

func InitServices(userRepo repository.UserRepo, forumRepo repository.ForumRepo, threadRepo repository.ThreadRepo, postRepo repository.PostRepo, voteRepo repository.VoteRepo, statusRepo repository.StatusRepo) (delivery.UserService, delivery.ForumService, delivery.ThreadService, delivery.PostService, delivery.VoteService, delivery.StatusService) {
	userService := delivery.MakeUserService(userRepo)
	forumService := delivery.MakeForumService(forumRepo)
	threadService := delivery.MakeThreadService(threadRepo)
	postService := delivery.MakePostService(postRepo)
	voteService := delivery.MakeVoteService(voteRepo)
	serviceService := delivery.MakeServiceService(statusRepo)

	return userService, forumService, threadService, postService, voteService, serviceService
}

func main() {
	db := InitDb()
	defer db.Close()
	userRepo, forumRepo, threadRepo, postRepo, voteRepo, statusRepo := InitRepos(db)
	userService, forumService, threadService, postService, voteService, serviceService := InitServices(userRepo, forumRepo, threadRepo, postRepo, voteRepo, statusRepo)
	apiRouter := router.New()

	apiRouter.POST("/api/forum/create", forumService.Create)
	apiRouter.GET("/api/forum/{slug}/details", forumService.Get)
	apiRouter.POST("/api/forum/{slug}/create", threadService.Create)
	apiRouter.GET("/api/forum/{slug}/users", forumService.GetUsers)
	apiRouter.GET("/api/forum/{slug}/threads", forumService.GetThreads)

	apiRouter.POST("/api/thread/{slug_or_id}/details", threadService.Update)
	apiRouter.GET("/api/thread/{slug_or_id}/posts", threadService.GetPosts)
	apiRouter.POST("/api/thread/{slug_or_id}/vote", voteService.Create)
	apiRouter.POST("/api/thread/{slug_or_id}/create", postService.Create)
	apiRouter.GET("/api/thread/{slug_or_id}/details", threadService.Get)

	apiRouter.GET("/api/post/{id}/details", postService.Get)
	apiRouter.POST("/api/post/{id}/details", postService.Update)

	apiRouter.POST("/api/user/{nickname}/create", userService.Create)
	apiRouter.GET("/api/user/{nickname}/profile", userService.Get)
	apiRouter.POST("/api/user/{nickname}/profile", userService.Update)

	apiRouter.GET("/api/service/status", serviceService.GetStatus)
	apiRouter.POST("/api/service/clear", serviceService.Clear)

	cont := context.Background()

	log.Println("Запуск сервера...")
	err := fasthttp.ListenAndServe(
		"0.0.0.0:5000",
		func(context *fasthttp.RequestCtx) {
			context.SetUserValue("cont", cont)
			apiRouter.Handler(context)
		},
	)

	if err != nil {
		fmt.Printf("error on listening: %v", err)
	}
}
