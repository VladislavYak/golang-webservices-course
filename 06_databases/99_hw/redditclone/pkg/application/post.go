package application

import (
	"context"
	"errors"
	"fmt"

	postP "github.com/VladislavYak/redditclone/pkg/domain/post"
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
	returnedPost, err := p.repo.AddPost(ctx, Post)

	// как будто тут я хочу делать Upvote для своего поста. но тогда сюда надо параметры передатьва...

	fmt.Println("returnedPost", returnedPost)

	return returnedPost, err
}

func (p *PostImpl) Delete(ctx context.Context, s string) (*postP.Post, error) {

	returnedPost, err := p.repo.DeletePost(ctx, s)

	fmt.Println("returnedPost", returnedPost)

	return returnedPost, err

}

// yakovlev: тту по идее поитенр не нужен, но его и из интерфейса надо выпилить бы
func (p *PostImpl) GetAll(ctx context.Context) ([]*postP.Post, error) {

	returnedPost, err := p.repo.GetAllPosts(ctx)

	return returnedPost, err

}

func (p *PostImpl) GetByID(ctx context.Context, s string) (*postP.Post, error) {

	returnedPost, err := p.repo.GetPostByID(ctx, s)

	return returnedPost, err

}

// yakovlev: тту по идее поитенр не нужен, но его и из интерфейса надо выпилить бы
func (p *PostImpl) GetPostsByCategoryName(ctx context.Context, s string) ([]*postP.Post, error) {

	returnedPost, err := p.repo.GetPostsByCategoryName(ctx, s)

	return returnedPost, err

}

func (p *PostImpl) GetByUsername(ctx context.Context, s string) ([]*postP.Post, error) {

	returnedPost, err := p.repo.GetPostsByUsername(ctx, s)

	return returnedPost, err
}

func (p *PostImpl) Upvote(ctx context.Context, PostId string) (*postP.Post, error) {

	UserID, ok := ctx.Value("UserID").(string)
	if !ok {
		return nil, errors.New("cannot cast userID to string")
	}

	returnedPost, err := p.repo.Upvote(ctx, PostId, UserID)

	return returnedPost, err
}

func (p *PostImpl) Downvote(ctx context.Context, PostId string) (*postP.Post, error) {
	UserID, ok := ctx.Value("UserID").(string)
	if !ok {
		return nil, errors.New("cannot cast userID to string")
	}

	returnedPost, err := p.repo.Downvote(ctx, PostId, UserID)

	return returnedPost, err
}

func (p *PostImpl) Unvote(ctx context.Context, PostId string) (*postP.Post, error) {
	UserID, ok := ctx.Value("UserID").(string)
	if !ok {
		return nil, errors.New("cannot cast userID to string")
	}

	returnedPost, err := p.repo.Unvote(ctx, PostId, UserID)

	return returnedPost, err
}
