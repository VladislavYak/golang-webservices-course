package handlers

import (
	"fmt"
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/post"
	"github.com/VladislavYak/redditclone/pkg/user"
	"github.com/golang-jwt/jwt/v5"
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

func (ph *PostHandler) GetPostByID(c echo.Context) error {
	id := c.Param("id")

	if err := ph.Repo.UpdatePostViews(id); err != nil {
		return err
	}

	post, err := ph.Repo.GetPostByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)

}

func (ph *PostHandler) GetPostByUsername(c echo.Context) error {
	username := c.Param("username")

	post, err := ph.Repo.GetPostsByUsername(username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)

}

func (ph *PostHandler) PostPost(c echo.Context) error {
	// yakovlev: корнер кейс, когда текст - есть только текст, но нет урл
	// когда урл, нет текст
	us := c.Get("user").(*jwt.Token)
	claims := us.Claims.(*JwtCustomClaims)

	pp := &PostParams{}
	if err := c.Bind(pp); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	Post := post.NewPost(pp.Category, pp.Type, pp.Url, pp.Text, pp.Title, *user.NewUser(claims.Name).WithID(claims.Id))

	postReturned, err := ph.Repo.AddPost(Post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot add post")
	}

	return c.JSON(http.StatusCreated, postReturned)
}

// yakovlev: use repo here
func (ph *PostHandler) DeletePost(c echo.Context) error {

	id := c.Param("id")
	// idInt, err := strconv.Atoi(id)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "got invalid id")
	// }

	for i, value := range ph.Repo.Data {
		if value.Id == id {
			ph.Repo.Data = append(ph.Repo.Data[:i], ph.Repo.Data[i+1:]...)
		}
		return c.JSON(http.StatusOK, value)
	}

	return echo.NewHTTPError(http.StatusNotFound, "this id doesnot exist")
}

func (ph *PostHandler) AddComment(c echo.Context) error {
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

	Comment := post.NewComment(*user.NewUser(claims.Name).WithID(claims.Id), body.Comment)

	returnedPost, err := ph.Repo.AddComment(id, Comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusCreated, returnedPost)
}

func (ph *PostHandler) DeleteComment(c echo.Context) error {
	id := c.Param("id")
	commentId := c.Param("commentId")

	ph.Repo.Mutex.Lock()
	defer ph.Repo.Mutex.Unlock()

	for i, post := range ph.Repo.Data {
		if post.Id == id {

			for j, comment := range post.Comments {
				if comment.Id == commentId {
					post.Comments = append(post.Comments[:j], post.Comments[j+1:]...)
					ph.Repo.Data[i] = post
					return c.JSON(http.StatusOK, post)
				}

			}

		}

	}
	return echo.NewHTTPError(http.StatusNotFound, "this id doesnot exist")
}
