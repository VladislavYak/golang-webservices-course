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
