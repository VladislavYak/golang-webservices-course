package post

import (
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
)

// need adding more fields
type Post struct {
	Id               string            `json:"id"`
	Category         string            `json:"category"`
	Type             string            `json:"type"`
	Url              string            `json:"url,omitempty"`
	Text             string            `json:"text,omitempty"`
	Title            string            `json:"title"`
	Votes            []Vote            `json:"votes"`
	Comments         []comment.Comment `json:"comments"`
	Created          time.Time         `json:"created"`
	UpvotePercentage int               `json:"upvotePercentage"`

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
		Comments: []comment.Comment{},
		Votes:    []Vote{},
	}
}

func (p *Post) WithId(id string) *Post {
	p.Id = id
	return p
}
