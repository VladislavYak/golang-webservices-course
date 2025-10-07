package comment

import (
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/user"
)

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

type CommentRepository interface {
	AddComment(string, *Comment) error
	DeleteComment(string, string) error
}
