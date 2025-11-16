package handlers

import (
	"errors"
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
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
	posts, err := ph.Implementation.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) GetPostsByCategoryName(c echo.Context) error {

	CategoryName := c.Param("CategoryName")

	posts, err := ph.Implementation.GetPostsByCategoryName(c.Request().Context(), CategoryName)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}

func (ph *PostHandler) GetPostByID(c echo.Context) error {
	id := c.Param("id")

	post, err := ph.Implementation.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, post)

}

func (ph *PostHandler) GetPostByUsername(c echo.Context) error {
	username := c.Param("username")

	post, err := ph.Implementation.GetPostsByUsername(c.Request().Context(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, post)

}

func (ph *PostHandler) PostPost(c echo.Context) error {

	ctx := getUserCtx(c)

	userID, ok := ctx.Value("UserID").(string)
	if !ok {
		return errors.New("cannot cast userID to string")
	}

	username, ok := ctx.Value("Username").(string)
	if !ok {
		return errors.New("cannot cast userID to string")
	}

	pp := &PostParams{}
	if err := c.Bind(pp); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	Post := post.NewPost(pp.Category, pp.Type, pp.Url, pp.Text, pp.Title, *user.NewUser(username).WithID(userID))

	postReturned, err := ph.Implementation.Create(ctx, Post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, postReturned)
}

func (ph *PostHandler) DeletePost(c echo.Context) error {

	id := c.Param("id")

	ctx := getUserCtx(c)

	userID, ok := ctx.Value("UserID").(string)
	if !ok {
		return errors.New("cannot cast userID to string")
	}

	deletedPost, err := ph.Implementation.Delete(ctx, id, userID)

	if errors.Is(err, post.DifferentPostOwnerError) {
		return echo.NewHTTPError(http.StatusForbidden, err)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, deletedPost)
}

func (ph *PostHandler) Upvote(c echo.Context) error {
	PostId := c.Param("id")

	ctx := getUserCtx(c)

	returnedPost, err := ph.Implementation.Upvote(ctx, PostId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, returnedPost)

}

func (ph *PostHandler) Downvote(c echo.Context) error {
	PostId := c.Param("id")

	ctx := getUserCtx(c)

	returnedPost, err := ph.Implementation.Downvote(ctx, PostId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, returnedPost)
}

func (ph *PostHandler) Unvote(c echo.Context) error {
	PostId := c.Param("id")

	ctx := getUserCtx(c)

	returnedPost, err := ph.Implementation.Unvote(ctx, PostId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, returnedPost)

}
