package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladislavYak/redditclone/pkg/post"
	"github.com/labstack/echo/v4"
)

type PostParams struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Url      string `json:"url"`
	Text     string `json:"text"`
	Title    string `json:"title"`
}

type PostHandler struct {
	Repo post.PostRepo
}

func (ph *PostHandler) GetPosts(c echo.Context) error {
	posts, err := ph.Repo.GetAllPosts()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) GetPostsByCategoryName(c echo.Context) error {
	CategoryName := c.Param("CategoryName")

	posts, err := ph.Repo.GetPostsByCategoryName(CategoryName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) PostPost(c echo.Context) error {
	pp := &PostParams{}
	if err := c.Bind(pp); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	fmt.Println("pp", pp)

	Post := &post.Post{Id: 0, Category: pp.Category, Type: pp.Type, Url: pp.Url, Text: pp.Text, Title: pp.Title}
	postReturned, err := ph.Repo.AddPost(Post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot add post")
	}

	return c.JSON(http.StatusCreated, postReturned)
}

func (ph *PostHandler) DeletePost(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "got invalid id")
	}

	for i, value := range ph.Repo.Data {
		if value.Id == idInt {
			ph.Repo.Data = append(ph.Repo.Data[:i], ph.Repo.Data[i+1:]...)
		}
		return c.JSON(http.StatusOK, value)
	}

	return echo.NewHTTPError(http.StatusNotFound, "this id doesnot exist")
}

func remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
