package ram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/VladislavYak/redditclone/pkg/domain/post"
)

var _ post.PostRepository = new(PostRepo)

type PostRepo struct {
	Data []*post.Post
	*sync.Mutex
	lastID    int
	commentID int
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		[]*post.Post{},
		&sync.Mutex{},
		0,
		0,
	}
}

func (pp *PostRepo) GetAllPosts(ctx context.Context) ([]*post.Post, error) {
	return pp.Data, nil
}

func (pp *PostRepo) GetPostsByCategoryName(ctx context.Context, CategoryName string) ([]*post.Post, error) {
	res := []*post.Post{}

	for _, post := range pp.Data {
		if post.Category == CategoryName {
			res = append(res, post)
		}
	}
	return res, nil

}

func (pp *PostRepo) GetPostByID(ctx context.Context, ID string) (*post.Post, error) {
	for _, post := range pp.Data {
		if post.Id == ID {
			return post, nil
		}
	}
	return nil, errors.New("not found")
}

func (pp *PostRepo) GetPostsByUsername(ctx context.Context, Username string) ([]*post.Post, error) {
	res := []*post.Post{}

	for _, post := range pp.Data {
		if post.Author.Username == Username {
			res = append(res, post)
		}
	}
	return res, nil

}

func (pp *PostRepo) UpdatePostViews(ctx context.Context, ID string) error {
	for _, Post := range pp.Data {
		if Post.Id == ID {
			Post.Views += 1
			return nil
		}
	}
	return errors.New("not found")
}

func (pp *PostRepo) AddPost(ctx context.Context, Post *post.Post) (*post.Post, error) {
	// pretty random handling mutexes
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	Post = Post.WithId(strconv.Itoa(pp.lastID))

	pp.lastID++

	pp.Data = append(pp.Data, Post)

	fmt.Println("my Posts", pp.Data)
	return Post, nil
}

func (pp *PostRepo) DeletePost(ctx context.Context, Id string) (*post.Post, error) {
	for i, value := range pp.Data {
		if value.Id == Id {
			pp.Data = append(pp.Data[:i], pp.Data[i+1:]...)
		}
		return value, nil
	}

	return nil, errors.New("this id doesnot exist")

}

// yakovlev: add proper error handling
func (pp *PostRepo) Upvote(ctx context.Context, PostId string, UserId string) (*post.Post, error) {
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for i, Post := range pp.Data {
		if Post.Id == PostId {
			for j, voteIter := range Post.Votes {
				if voteIter.User == UserId {

					pp.Data[i].Votes[j].WithVote(1)
					pp.Data[i].UpdateScore()
					return pp.Data[i], nil
				}
			}

			pp.Data[i].Votes = append(pp.Data[i].Votes, post.Vote{User: UserId, VoteScore: 1})
			// Post.Votes = append(Post.Votes, Vote{User: user_id, VoteScore: -1})

			pp.Data[i].UpdateScore()

			return pp.Data[i], nil
		}
	}

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepo) Downvote(ctx context.Context, id string, UserId string) (*post.Post, error) {
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for i, Post := range pp.Data {
		if Post.Id == id {
			for j, voteIter := range Post.Votes {
				if voteIter.User == UserId {

					pp.Data[i].Votes[j].WithVote(-1)
					pp.Data[i].UpdateScore()
					return pp.Data[i], nil
				}
			}

			pp.Data[i].Votes = append(pp.Data[i].Votes, post.Vote{User: UserId, VoteScore: -1})
			// Post.Votes = append(Post.Votes, Vote{User: user_id, VoteScore: -1})

			pp.Data[i].UpdateScore()

			return pp.Data[i], nil
		}
	}

	return nil, errors.New("this id doesnot exist")
}

func (pp *PostRepo) Unvote(ctx context.Context, id string, UserId string) (*post.Post, error) {
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for i, Post := range pp.Data {
		if Post.Id == id {
			for j, voteIter := range Post.Votes {
				if voteIter.User == UserId {

					pp.Data[i].Votes = append(pp.Data[i].Votes[:j], pp.Data[i].Votes[j+1:]...)
					pp.Data[i].UpdateScore()
					return pp.Data[i], nil
				}
			}

			return nil, errors.New("cannot find user for specified post")
		}
	}

	return nil, errors.New("this id doesnot exist")
}
