package handlers

import (
	"fmt"
	"net/http"

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

// currently
// {
//     "id": "0",
//     "category": "music",
//     "type": "text",
//     "url": "testtesttet",
//     "text": "tetetetet",
//     "title": "tettetet",
//     "votes": null,
//     "created": "",
//     "upvotePercentage": 0,
//     "score": 0,
//     "views": 0,
//     "Author": {
//         "username": "",
//         "password": "",
//         "id": 0
//     }
// }

// golden state
//
//	{
//	    "score": 1,
//	    "views": 4,
//	    "type": "text",
//	    "title": "ddsadsa",
//	    "author": {
//	        "username": "dasdasddas",
//	        "id": "656e23ad1d06de00132f7e20"
//	    },
//	    "category": "music",
//	    "text": "dasdsadas",
//	    "votes": [
//	        {
//	            "user": "656e23ad1d06de00132f7e20",
//	            "vote": 1
//	        }
//	    ],
//	    "comments": [],
//	    "created": "2023-12-04T19:23:02.906Z",
//	    "upvotePercentage": 100,
//	    "id": "656e27161d06de00132f7e25"
//	}
//
// Определяем структуру для JSON-ответа
type Response struct {
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Type             string    `json:"type"`
	Title            string    `json:"title"`
	Author           Author    `json:"author"`
	Category         string    `json:"category"`
	Text             string    `json:"text"`
	Votes            []Vote    ` json:"votes"`
	Comments         []Comment `json:"comments"`
	Created          string    `json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`
	ID               string    `json:"id"`
}

type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type Vote struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type Comment struct {
	// Добавьте поля для комментариев, если необходимо
}

func (ph *PostHandler) GetPostByID(c echo.Context) error {
	// id := c.Param("id")

	// idConv, err := strconv.Atoi(id)
	// if err != nil {
	// 	return errors.New("got invalid id")
	// }

	// post, err := ph.Repo.GetPostByID(id)
	// if err != nil {
	// 	return err
	// }
	// Создаем фиксированный ответ
	response := Response{
		Score: 1,
		Views: 4,
		Type:  "text",
		Title: "ddsadsa",
		Author: Author{
			Username: "dasdasddas",
			ID:       "656e23ad1d06de00132f7e20",
		},
		Category: "music",
		Text:     "dasdsadas",
		Votes: []Vote{
			{
				User: "656e23ad1d06de00132f7e20",
				Vote: 1,
			},
		},
		Comments:         []Comment{},
		Created:          "2023-12-04T19:23:02.906Z",
		UpvotePercentage: 100,
		ID:               "656e27161d06de00132f7e25",
	}
	return c.JSON(http.StatusOK, response)

}

func (ph *PostHandler) PostPost(c echo.Context) error {
	pp := &PostParams{}
	if err := c.Bind(pp); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	fmt.Println("pp", pp)

	Post := &post.Post{Category: pp.Category, Type: pp.Type, Url: pp.Url, Text: pp.Text, Title: pp.Title}
	postReturned, err := ph.Repo.AddPost(Post)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot add post")
	}
	return c.JSON(http.StatusCreated, postReturned)
}

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
