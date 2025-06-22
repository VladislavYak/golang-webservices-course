package post

import "github.com/VladislavYak/redditclone/pkg/user"

type Vote struct {
	User int `json:"id"`
	Vote int `json:"vote"`
}

type Comment struct {
	Created string    `json:"created"`
	Author  user.User `json:"author"`
	Body    string    `json:"body"`
	Id      string    `json:"id"`
}

// need adding more fields
type Post struct {
	Id               string    `json:"id"`
	Category         string    `json:"category"`
	Type             string    `json:"type"`
	Url              string    `json:"url"`
	Text             string    `json:"text"`
	Title            string    `json:"title"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"`
	Created          string    `json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`

	Score int `json:"score"`
	Views int `json:"views"`

	// user should be deleted by password
	Author user.User
}
