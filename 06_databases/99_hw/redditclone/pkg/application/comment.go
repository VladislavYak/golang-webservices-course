package application

import (
	"context"

	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	commentC "github.com/VladislavYak/redditclone/pkg/domain/comment"
	postP "github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/go-faster/errors"
)

type CommentInterface interface {
	AddComment(context.Context, string, *comment.Comment) (*postP.Post, error)
	DeleteComment(context.Context, string, string) (*postP.Post, error)
}

var _ CommentInterface = new(CommentImpl)

type CommentImpl struct {
	PostRepo    postP.PostRepository
	CommentRepo commentC.CommentRepository
}

func NewCommentImpl(repoP postP.PostRepository, repoC commentC.CommentRepository) *CommentImpl {
	return &CommentImpl{PostRepo: repoP, CommentRepo: repoC}
}

// yakovlev: сейчас можно пустой коммент оставить - а это плохо)
func (ci *CommentImpl) AddComment(ctx context.Context, id string, Comment *comment.Comment) (*postP.Post, error) {
	const op = "AddComment"
	err := ci.CommentRepo.AddComment(ctx, id, Comment)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	returnedPost, err := ci.PostRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, nil
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6InRlc3R0ZXN0MiIsInVzZXJfaWQiOiI0IiwiZXhwIjoxNzYzMzkwODI1LCJpYXQiOjE3NjMzODk5MjV9.K0nH0a_ZiRNqtZM2Al_TVc_dn5OmrKPYN_z47lbp8FI

// сейчас любой чел может удалить любой коммент? - da
func (ci *CommentImpl) DeleteComment(ctx context.Context, PostId string, CommentId string) (*postP.Post, error) {
	const op = "DeleteComment"
	err := ci.CommentRepo.DeleteComment(ctx, PostId, CommentId)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	returnedPost, err := ci.PostRepo.GetPostByID(ctx, PostId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, nil
}
