package post

import (
	"fmt"
	"time"

	"github.com/VladislavYak/redditclone/pkg/user"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Vote struct {
	User      string `json:"id"`
	VoteScore int    `json:"vote"`
}

func (v *Vote) WithVote(value int) *Vote {
	v.VoteScore = value

	return v
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
	MongoId          MongoObjectID `bson:"_id,omitempty" json:"-"`
	Id               string        `json:"id"`
	Category         string        `json:"category"`
	Type             string        `json:"type"`
	Url              string        `json:"url,omitempty"`
	Text             string        `json:"text,omitempty"`
	Title            string        `json:"title"`
	Votes            []Vote        `json:"votes"`
	Comments         []Comment     `json:"comments"`
	Created          time.Time     `json:"created"`
	UpvotePercentage int           `json:"upvotePercentage"`

	Score int `json:"score"`
	Views int `json:"views"`

	Author user.User `json:"author"`
}

type MongoObjectID string

func (moi *MongoObjectID) UnmarshalBSONValue(t byte, b []byte) error {
	fmt.Println("im inside UnmarshalBSONValue...")

	if bson.Type(t) == bson.TypeObjectID {

		var oid bson.ObjectID
		copy(oid[:], b)

		fmt.Println("oid", oid)

		*moi = MongoObjectID(oid.Hex())

		fmt.Println("oid.Hex()", oid.Hex())
	}
	return nil
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

func (p *Post) UpdateScore() *Post {
	updatedScore := 0
	for _, vote := range p.Votes {
		updatedScore += vote.VoteScore
	}

	p.Score = updatedScore
	return p
}
