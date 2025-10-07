package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	Implementation application.CommentInterface
}

func (ch *CommentHandler) AddComment(c echo.Context) error {
	us := c.Get("user").(*jwt.Token)
	claims := us.Claims.(*JwtCustomClaims)

	id := c.Param("id")

	var body struct {
		Comment string `json:"comment"`
	}

	fmt.Println("body", body)

	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	Comment := comment.NewComment(*user.NewUser(claims.Name).WithID(claims.Id), body.Comment)

	returnedPost, err := ch.Implementation.AddComment(context.TODO(), id, Comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusCreated, returnedPost)
}

func (ch *CommentHandler) DeleteComment(c echo.Context) error {
	id := c.Param("id")
	CommentId := c.Param("commentId")

	returnedPost, err := ch.Implementation.DeleteComment(context.TODO(), id, CommentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return echo.NewHTTPError(http.StatusCreated, returnedPost)
}
