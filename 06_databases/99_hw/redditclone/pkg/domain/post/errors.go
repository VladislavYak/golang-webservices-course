package post

import "github.com/go-faster/errors"

var (
	DifferentPostOwnerError = errors.New("cannot delete post: belongs to another user")
)
