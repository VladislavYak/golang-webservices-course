package post

import "context"

type PostRepository interface {
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPostsByCategoryName(ctx context.Context, CategoryName string) ([]*Post, error)
	GetPostByID(ctx context.Context, ID string) (*Post, error)
	GetPostsByUsername(ctx context.Context, Username string) ([]*Post, error)
	AddPost(ctx context.Context, Post *Post) (*Post, error)
	DeletePost(ctx context.Context, Id string) (*Post, error)
	Vote(ctx context.Context, PostId string, vote int) (*Post, error)
	// Downvote(ctx context.Context, Id string) (*Post, error)
	Unvote(ctx context.Context, Id string) (*Post, error)
	UpdateScore(ct context.Context, Id string) (*Post, error)
}
