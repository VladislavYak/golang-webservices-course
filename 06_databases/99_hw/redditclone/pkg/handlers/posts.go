package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostParams struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Url      string `json:"url,omitempty"`
	Text     string `json:"text,omitempty"`
	Title    string `json:"title"`
}

type PostHandler struct {
	Implementation application.PostInterface
}

func (ph *PostHandler) GetPosts(c echo.Context) error {
	posts, err := ph.Implementation.GetAll(context.TODO())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) GetPostsByCategoryName(c echo.Context) error {

	CategoryName := c.Param("CategoryName")

	posts, err := ph.Implementation.GetPostsByCategoryName(context.TODO(), CategoryName)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get posts")
	}

	return c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) GetPostByID(c echo.Context) error {
	id := c.Param("id")

	// if err := ph.Repo.UpdatePostViews(id); err != nil {
	// 	return err
	// }

	post, err := ph.Implementation.GetByID(context.TODO(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)

}

func (ph *PostHandler) GetPostByUsername(c echo.Context) error {
	username := c.Param("username")

	post, err := ph.Implementation.GetByUsername(context.TODO(), username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)

}

func (ph *PostHandler) PostPost(c echo.Context) error {

	us := c.Get("user").(*jwt.Token)

	// yakovlev: do this at middleware somehow
	claims := us.Claims.(*JwtCustomClaims)

	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "username", claims.Name)
	ctx = context.WithValue(ctx, "id", claims.Id)

	pp := &PostParams{}
	if err := c.Bind(pp); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	fmt.Println("username (handler)", ctx.Value("username"))
	fmt.Println("UserID (handler)", ctx.Value("id"))
	Post := post.NewPost(pp.Category, pp.Type, pp.Url, pp.Text, pp.Title, *user.NewUser(claims.Name).WithID(claims.Id))

	postReturned, err := ph.Implementation.Create(ctx, Post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot add post")
	}

	return c.JSON(http.StatusCreated, postReturned)
}

func (ph *PostHandler) DeletePost(c echo.Context) error {

	id := c.Param("id")

	deletedPost, err := ph.Implementation.Delete(context.TODO(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, deletedPost)
}

func (ph *PostHandler) Upvote(c echo.Context) error {
	PostId := c.Param("id")

	us := c.Get("user").(*jwt.Token)
	claims := us.Claims.(*JwtCustomClaims)

	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "Username", claims.Name)
	ctx = context.WithValue(ctx, "UserID", claims.Id)

	returnedPost, err := ph.Implementation.Upvote(ctx, PostId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK, returnedPost)

}

func (ph *PostHandler) Downvote(c echo.Context) error {
	PostId := c.Param("id")

	us := c.Get("user").(*jwt.Token)
	claims := us.Claims.(*JwtCustomClaims)

	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "Username", claims.Name)
	ctx = context.WithValue(ctx, "UserID", claims.Id)

	returnedPost, err := ph.Implementation.Downvote(ctx, PostId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK, returnedPost)
}

func (ph *PostHandler) Unvote(c echo.Context) error {
	PostId := c.Param("id")

	us := c.Get("user").(*jwt.Token)
	claims := us.Claims.(*JwtCustomClaims)

	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "Username", claims.Name)
	ctx = context.WithValue(ctx, "UserID", claims.Id)

	returnedPost, err := ph.Implementation.Unvote(ctx, PostId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK, returnedPost)

}
