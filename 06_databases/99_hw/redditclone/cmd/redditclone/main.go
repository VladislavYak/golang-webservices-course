package main

import (
	"github.com/VladislavYak/redditclone/pkg/handlers"
	// "github.com/VladislavYak/redditclone/pkg/post"
	"github.com/VladislavYak/redditclone/pkg/application"

	"github.com/VladislavYak/redditclone/pkg/infrastructure/mongodb"
	"github.com/VladislavYak/redditclone/pkg/infrastructure/ram"
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

	client, _ := mongodb.NewMongoClient(cfg)
	PostRepo := mongodb.NewPostRepoMongo(client, "testing", "posts")
	// PostRepo := ram.NewPostRepo()

	// CommentRepo := ram.NewCommentRepo()
	CommentRepo := mongodb.NewCommentRepoMongo(client, "testing", "posts")
	UserRepo := ram.NewUserRepo()

	PostImpl := application.NewPostImpl(PostRepo)
	CommentImpl := application.NewCommentImpl(PostRepo, CommentRepo)

	postHandler := handlers.PostHandler{Implementation: PostImpl}
	commentHandler := handlers.CommentHandler{Implementation: CommentImpl}

	registerHandler := handlers.RegisterHandler{UserRepo: *UserRepo}
	loginHandler := handlers.LoginHandler{UserRepo: *UserRepo}

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Logger.Debug()

	g := e.Group(api)

	// Отдаем статические файлы
	e.Static("/static", "static")
	// e.Static("/", "static/html") // Отдаем HTML файлы по корневому маршруту
	e.File("/", "static/html/index.html")

	g.POST("/register", registerHandler.Register)
	g.POST("/login", loginHandler.Login)
	g.GET("/posts", postHandler.GetPosts)
	g.GET("/posts/:CategoryName", postHandler.GetPostsByCategoryName)
	g.GET("/post/:id", postHandler.GetPostByID)

	g.GET("/user/:username", postHandler.GetPostByUsername)

	{
		config := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(handlers.JwtCustomClaims)
			},
			SigningKey: []byte("secret"),
		}
		g.Use(echojwt.WithConfig(config))

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
