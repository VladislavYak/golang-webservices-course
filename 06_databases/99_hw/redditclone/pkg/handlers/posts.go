package handlers

import (
	"context"
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
	claims := us.Claims.(*JwtCustomClaims)

	pp := &PostParams{}
	if err := c.Bind(pp); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	Post := post.NewPost(pp.Category, pp.Type, pp.Url, pp.Text, pp.Title, *user.NewUser(claims.Name).WithID(claims.Id))

	postReturned, err := ph.Implementation.Create(context.TODO(), Post)
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

// func (ph *PostHandler) AddComment(c echo.Context) error {
// 	us := c.Get("user").(*jwt.Token)
// 	claims := us.Claims.(*JwtCustomClaims)

// 	id := c.Param("id")

// 	var body struct {
// 		Comment string `json:"comment"`
// 	}

// 	fmt.Println("body", body)

// 	if err := c.Bind(&body); err != nil {
// 		return c.String(http.StatusBadRequest, err.Error())
// 	}

// 	Comment := post.NewComment(*user.NewUser(claims.Name).WithID(claims.Id), body.Comment)

// 	returnedPost, err := ph.Repo.AddComment(id, Comment)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}

// 	return echo.NewHTTPError(http.StatusCreated, returnedPost)
// }

// func (ph *PostHandler) DeleteComment(c echo.Context) error {
// 	id := c.Param("id")
// 	commentId := c.Param("commentId")

// 	returnedPost, err := ph.Repo.DeleteComment(id, commentId)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}
// 	return echo.NewHTTPError(http.StatusCreated, returnedPost)
// }

// func (ph *PostHandler) Upvote(c echo.Context) error {
// 	id := c.Param("id")
// 	us := c.Get("user").(*jwt.Token)
// 	claims := us.Claims.(*JwtCustomClaims)

// 	returnedPost, err := ph.Repo.Upvote(id, claims.Id)

// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}

// 	return echo.NewHTTPError(http.StatusOK, returnedPost)

// }

// func (ph *PostHandler) Downvote(c echo.Context) error {
// 	id := c.Param("id")
// 	us := c.Get("user").(*jwt.Token)
// 	claims := us.Claims.(*JwtCustomClaims)

// 	returnedPost, err := ph.Repo.Downvote(id, claims.Id)

// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}

// 	return echo.NewHTTPError(http.StatusOK, returnedPost)
// }

// func (ph *PostHandler) Unvote(c echo.Context) error {
// 	id := c.Param("id")
// 	us := c.Get("user").(*jwt.Token)
// 	claims := us.Claims.(*JwtCustomClaims)

// 	returnedPost, err := ph.Repo.Unvote(id, claims.Id)

// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err)
// 	}

// 	return echo.NewHTTPError(http.StatusOK, returnedPost)

// }
