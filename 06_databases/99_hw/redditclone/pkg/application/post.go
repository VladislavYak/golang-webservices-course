package application

import (
	"context"
	"fmt"

	postP "github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/pkg/errors"
)

// https://gitlab.com/slax0rr/go-beer-api/-/blob/master/application/beer.go?ref_type=heads#L25
// тут я вижу, что определ интерфес. Мне не очень понятно. Зачем он нужен

type PostInterface interface {
	Create(context.Context, *postP.Post) (*postP.Post, error)
	Delete(context.Context, string) (*postP.Post, error)
	GetAll(context.Context) ([]*postP.Post, error)
	GetByID(context.Context, string) (*postP.Post, error)
	GetPostsByCategoryName(context.Context, string) ([]*postP.Post, error)
	GetByUsername(context.Context, string) ([]*postP.Post, error)
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

	fmt.Println("inside create")
	returnedPost, err := p.repo.AddPost(ctx, Post)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	fmt.Println("returnedPost.Id", returnedPost.Id)

	returnedPost, err = p.repo.Upvote(ctx, returnedPost.Id)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	fmt.Println("returnedPost", returnedPost)

	return returnedPost, err
}

// yakovlev: тут возвращать delete надо
func (p *PostImpl) Delete(ctx context.Context, s string) (*postP.Post, error) {
	const op = "Delete"
	returnedPost, err := p.repo.DeletePost(ctx, s)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	fmt.Println("returnedPost", returnedPost)

	return returnedPost, err

}

// yakovlev: тту по идее поитенр не нужен, но его и из интерфейса надо выпилить бы
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

// yakovlev: тту по идее поитенр не нужен, но его и из интерфейса надо выпилить бы
func (p *PostImpl) GetPostsByCategoryName(ctx context.Context, s string) ([]*postP.Post, error) {
	const op = "GetPostsByCategoryName"

	returnedPost, err := p.repo.GetPostsByCategoryName(ctx, s)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err

}

func (p *PostImpl) GetByUsername(ctx context.Context, s string) ([]*postP.Post, error) {
	const op = "GetByUsername"
	returnedPost, err := p.repo.GetPostsByUsername(ctx, s)

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return returnedPost, err
}

// fix updatescore
func (p *PostImpl) Upvote(ctx context.Context, PostId string) (*postP.Post, error) {
	const op = "Upvote"

	returnedPost, err := p.repo.Upvote(ctx, PostId)

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
func (p *PostImpl) Downvote(ctx context.Context, PostId string) (*postP.Post, error) {
	const op = "Downvote"
	returnedPost, err := p.repo.Downvote(ctx, PostId)

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
