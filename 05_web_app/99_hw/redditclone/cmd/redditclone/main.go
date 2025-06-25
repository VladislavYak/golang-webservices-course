package main

import (
	"github.com/VladislavYak/redditclone/pkg/handlers"
	"github.com/VladislavYak/redditclone/pkg/post"
	"github.com/VladislavYak/redditclone/pkg/user"
	jwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	api = "/api"
)

func main() {
	PostRepo := post.NewPostRepo()
	UserRepo := user.NewUserRepo()

	postHandler := handlers.PostHandler{Repo: *PostRepo}
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
		g.POST("/post/:id", postHandler.AddComment)
		g.DELETE("/post/:id/:commentId", postHandler.DeleteComment)
	}

	e.Logger.Fatal(e.Start(":1323"))

}
