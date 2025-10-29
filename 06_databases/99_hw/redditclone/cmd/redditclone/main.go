package main

import (
	"fmt"

	"github.com/VladislavYak/redditclone/pkg/handlers"
	// "github.com/VladislavYak/redditclone/pkg/post"
	"github.com/VladislavYak/redditclone/pkg/application"

	"github.com/VladislavYak/redditclone/pkg/domain/auth"
	customMiddleware "github.com/VladislavYak/redditclone/pkg/handlers/middleware"
	"github.com/VladislavYak/redditclone/pkg/infrastructure/mongodb"
	"github.com/VladislavYak/redditclone/pkg/infrastructure/postgres"
	jwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// yakovlev: add proper error handling
// yakovlev: add proper logging
// yakovlev: add proper concurrency handling using mutexes
const (
	api = "/api"
)

func main() {
	cfg := mongodb.Config{
		URI:        "mongodb://localhost:27017",
		Database:   "testing",
		TimeoutSec: 2,
	}

	// need to
	var JwtSecret = "secret"

	// better user DI somehow here

	client, _ := mongodb.NewMongoClient(cfg)
	pgPool, _ := postgres.NewPgPool()

	UserRepo := postgres.NewUserRepoPostgres(pgPool)
	// UserRepo := ram.NewUserRepo()
	PostRepo := mongodb.NewPostRepoMongo(client, "testing", "posts")
	// PostRepo := ram.NewPostRepo()
	AuthRepo := postgres.NewAuthRepoPostgres(pgPool)

	CommentRepo := mongodb.NewCommentRepoMongo(client, "testing", "posts")
	// CommentRepo := ram.NewCommentRepo()

	PostImpl := application.NewPostImpl(PostRepo)
	CommentImpl := application.NewCommentImpl(PostRepo, CommentRepo)
	UserImpl := application.NewUserImpl(UserRepo, AuthRepo, JwtSecret)
	AuthImpl := application.NewAuthImpl(AuthRepo)

	postHandler := handlers.PostHandler{Implementation: PostImpl}
	commentHandler := handlers.CommentHandler{Implementation: CommentImpl}

	userHandler := handlers.UserHandler{Impl: UserImpl}
	// loginHandler := handlers.LoginHandler{UserRepo: *UserRepo}

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Logger.Debug()

	g := e.Group(api)

	// Отдаем статические файлы
	e.Static("/static", "static")
	// e.Static("/", "static/html") // Отдаем HTML файлы по корневому маршруту
	e.File("/", "static/html/index.html")

	g.POST("/register", userHandler.Register)
	g.POST("/login", userHandler.Login)
	g.GET("/posts", postHandler.GetPosts)
	g.GET("/posts/:CategoryName", postHandler.GetPostsByCategoryName)
	g.GET("/post/:id", postHandler.GetPostByID)

	g.GET("/user/:username", postHandler.GetPostByUsername)

	{
		config := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(auth.JwtCustomClaims)
			},
			SigningKey: []byte("secret"),
		}

		fmt.Println("config", config)

		g.Use(customMiddleware.CustomAuth(&config, AuthImpl))
		// в общем кажется, что надо откащываться от этой withConfig и писать свою мидлварь для авторизации где есть проверка на валидность токена в бд
		// basicAuthMiddleware := echojwt.WithConfig(config)

		// где-то тут, наверное, мне нужна мидллаварь, которая ходит в базу и проверяет валидность токена.
		// но еще я делаб одинаковые операции с Claims. Вот их бы тоже унести куда-то.
		// эти операции делать надо после проверки мидлвар.
		// g.Use(basicAuthMiddleware)

		g.POST("/posts", postHandler.PostPost)
		g.DELETE("/post/:id", postHandler.DeletePost)
		g.POST("/post/:id", commentHandler.AddComment)
		g.DELETE("/post/:id/:commentId", commentHandler.DeleteComment)

		g.GET("/post/:id/downvote", postHandler.Downvote)
		g.GET("/post/:id/upvote", postHandler.Upvote)
		g.GET("/post/:id/unvote", postHandler.Unvote)
	}

	e.Logger.Fatal(e.Start(":1323"))

}
