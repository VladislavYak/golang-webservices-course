package handlers

import (
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/domain/auth"
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
	claims := us.Claims.(*auth.JwtCustomClaims)

	id := c.Param("id")

	var body struct {
		Comment string `json:"comment"`
	}

	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	Comment := comment.NewComment(*user.NewUser(claims.Login).WithID(claims.UserID), body.Comment)

	returnedPost, err := ch.Implementation.AddComment(c.Request().Context(), id, Comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, returnedPost)
}

func (ch *CommentHandler) DeleteComment(c echo.Context) error {
	ctx := getUserCtx(c)

	id := c.Param("id")
	CommentId := c.Param("commentId")

	returnedPost, err := ch.Implementation.DeleteComment(ctx, id, CommentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusOK, returnedPost)
}
