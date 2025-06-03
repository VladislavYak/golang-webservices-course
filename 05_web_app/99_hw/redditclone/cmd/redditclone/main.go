package main

import (
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/handlers"
	"github.com/VladislavYak/redditclone/pkg/post"
	"github.com/labstack/echo/v4"
)

const (
	api = "/api"
)

func main() {
	PostRepo := post.NewPostRepo()

	postHandler := handlers.PostHandler{Repo: *PostRepo}

	e := echo.New()
	e.Logger.Debug()

	g := e.Group(api)

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		c.Redirect(http.StatusSeeOther, "/static/html/index.html")
		return nil
	})

	g.GET("/posts", postHandler.GetPosts)
	g.GET("/posts/:CategoryName", postHandler.GetPostsByCategoryName)

	g.POST("/posts", postHandler.PostPost)
	g.DELETE("/post/:id", postHandler.DeletePost)

	g.POST("/register", handlers.Register)

	e.Logger.Fatal(e.Start(":1323"))

}
