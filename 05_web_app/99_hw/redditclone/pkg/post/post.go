package post

import (
	"time"

	"github.com/VladislavYak/redditclone/pkg/user"
)

type Vote struct {
	User int `json:"id"`
	Vote int `json:"vote"`
}

type Comment struct {
	Created time.Time `json:"created"`
	Author  user.User `json:"author"`
	Body    string    `json:"body"`
	Id      string    `json:"id"`
}

func NewComment(Author user.User, Body string) *Comment {
	return &Comment{Created: time.Now().UTC(), Author: Author, Body: Body}
}

func (c *Comment) WithId(id string) *Comment {
	c.Id = id
	return c
}

// need adding more fields
type Post struct {
	Id               string    `json:"id"`
	Category         string    `json:"category"`
	Type             string    `json:"type"`
	Url              string    `json:"url,omitempty"`
	Text             string    `json:"text,omitempty"`
	Title            string    `json:"title"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"`
	Created          time.Time `json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`

	Score int `json:"score"`
	Views int `json:"views"`

	Author user.User `json:"author"`
}

func NewPost(category string, postType string, url string, text string, title string, author user.User) *Post {
	return &Post{
		Category: category,
		Type:     postType,
		Url:      url,
		Text:     text,
		Title:    title,
		Author:   author,
		Created:  time.Now().UTC(),
		Score:    0,
		Views:    0,
		Comments: []Comment{},
	}
}

func (p *Post) WithId(id string) *Post {
	p.Id = id
	return p
}
