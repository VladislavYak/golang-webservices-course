package post

import (
	"context"
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
	}
}

func (p *Post) WithId(id string) *Post {
	p.Id = id
	return p
}

func (p *Post) UpdateScore() *Post {
	updatedScore := 0
	for _, vote := range p.Votes {
		updatedScore += vote.VoteScore
	}

	p.Score = updatedScore
	return p
}

type PostRepository interface {
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostsByCategoryName(ctx context.Context, CategoryName string) ([]*Post, error)
	GetPostByID(ctx context.Context, ID string) (*Post, error)
	GetPostsByUsername(ctx context.Context, Username string) ([]*Post, error)
	AddPost(ctx context.Context, Post *Post) (*Post, error)
	DeletePost(ctx context.Context, Id string) (*Post, error)
	// Save(ctx context.Context, user *domain.User) error
}

type Vote struct {
	User      string `json:"id"`
	VoteScore int    `json:"vote"`
}

func (v *Vote) WithVote(value int) *Vote {
	v.VoteScore = value

	return v
}
