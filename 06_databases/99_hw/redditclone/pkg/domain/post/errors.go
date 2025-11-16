package post

import "github.com/go-faster/errors"

var (
	DifferentPostOwnerError = errors.New("cannot delete post: belongs to another user")
	PostNotFoundError       = errors.New("post not found")
	InvalidPostIdError      = errors.New("invalid post id")
)
