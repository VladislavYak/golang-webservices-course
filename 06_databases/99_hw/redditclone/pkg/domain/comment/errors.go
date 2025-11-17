package comment

import "errors"

var (
	CommentNotFoundError        = errors.New("comment not found")
	DifferentCommentWriterError = errors.New("cannot delete comment of another user")
)
