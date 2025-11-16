package comment

import "context"

type CommentRepository interface {
	AddComment(context.Context, string, *Comment) error
	DeleteComment(context.Context, string, string) error
}
