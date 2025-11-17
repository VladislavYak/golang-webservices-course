package application

import (
	"context"

	"github.com/VladislavYak/redditclone/pkg/domain/post"
	postP "github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/pkg/errors"
)

// https://gitlab.com/slax0rr/go-beer-api/-/blob/master/application/beer.go?ref_type=heads#L25
// тут я вижу, что определ интерфес. Мне не очень понятно. Зачем он нужен

type PostInterface interface {
	Create(context.Context, *postP.Post) (*postP.Post, error)
	Delete(context.Context, string, string) (*postP.Post, error)
	GetAll(context.Context) ([]*postP.Post, error)
	GetByID(context.Context, string) (*postP.Post, error)
	GetPostsByCategoryName(context.Context, string) ([]*postP.Post, error)
	GetPostsByUsername(context.Context, string) ([]*postP.Post, error)
	Upvote(context.Context, string) (*postP.Post, error)
	Downvote(context.Context, string) (*postP.Post, error)
	Unvote(context.Context, string) (*postP.Post, error)
}

type PostImpl struct {
	repo postP.PostRepository
}

func NewPostImpl(repo postP.PostRepository) *PostImpl {
	return &PostImpl{repo: repo}
}

// Compile-time check if PostImpl implements PostInterface
var _ PostInterface = new(PostImpl)

func (p *PostImpl) Create(ctx context.Context, Post *postP.Post) (*postP.Post, error) {
	const op = "Create"

	returnedPost, err := p.repo.AddPost(ctx, Post)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	returnedPost, err = p.repo.Vote(ctx, returnedPost.Id, 1)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	returnedPost, err = p.repo.UpdateScore(ctx, returnedPost.Id)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

func (p *PostImpl) Delete(ctx context.Context, id string, userId string) (*postP.Post, error) {
	const op = "Delete"

	// yakovlev: in highload this could be raced condition
	ppost, err := p.repo.GetPostByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	if ppost.Author.UserID != userId {
		return nil, post.DifferentPostOwnerError
	}

	returnedPost, err := p.repo.DeletePost(ctx, id)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err

}

func (p *PostImpl) GetAll(ctx context.Context) ([]*postP.Post, error) {
	const op = "GetAll"

	returnedPost, err := p.repo.GetAllPosts(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

func (p *PostImpl) GetByID(ctx context.Context, s string) (*postP.Post, error) {
	const op = "GetByID"

	returnedPost, err := p.repo.GetPostByID(ctx, s)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

func (p *PostImpl) GetPostsByCategoryName(ctx context.Context, s string) ([]*postP.Post, error) {
	const op = "GetPostsByCategoryName"

	returnedPost, err := p.repo.GetPostsByCategoryName(ctx, s)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

func (p *PostImpl) GetPostsByUsername(ctx context.Context, s string) ([]*postP.Post, error) {
	const op = "GetPostsByUsername"
	returnedPost, err := p.repo.GetPostsByUsername(ctx, s)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

func (p *PostImpl) Upvote(ctx context.Context, PostId string) (*postP.Post, error) {
	const (
		op       = "Upvote"
		voteSign = 1
	)

	returnedPost, err := p.repo.Vote(ctx, PostId, voteSign)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	returnedPost, err = p.repo.UpdateScore(ctx, PostId)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

func (p *PostImpl) Downvote(ctx context.Context, PostId string) (*postP.Post, error) {
	const (
		op       = "Downvote"
		voteSign = -1
	)

	returnedPost, err := p.repo.Vote(ctx, PostId, voteSign)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	returnedPost, err = p.repo.UpdateScore(ctx, PostId)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

// fix updatescore
func (p *PostImpl) Unvote(ctx context.Context, PostId string) (*postP.Post, error) {
	const op = "Unvote"
	returnedPost, err := p.repo.Unvote(ctx, PostId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	returnedPost, err = p.repo.UpdateScore(ctx, PostId)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}
