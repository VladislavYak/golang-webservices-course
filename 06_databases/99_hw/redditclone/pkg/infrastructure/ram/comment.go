package ram

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
)

var _ comment.CommentRepository = new(CommentRepo)

type CommentRepo struct {
	Data []*post.Post
	*sync.Mutex
	lastID    int
	commentID int
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{}
}

func (pp *CommentRepo) AddComment(ctx context.Context, Id string, comment *comment.Comment) error {
	// add more mutexes handling
	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()

	for _, Post := range pp.Data {
		if Post.Id == Id {
			Post.Comments = append(Post.Comments, *comment.WithId(strconv.Itoa(pp.commentID)))

			pp.commentID++
			return nil
		}
	}

	return errors.New("post not found")
}

func (pp *CommentRepo) DeleteComment(ctx context.Context, id string, commentId string) error {

	pp.Mutex.Lock()
	defer pp.Mutex.Unlock()
	for i, post := range pp.Data {
		if post.Id == id {

			for j, comment := range post.Comments {
				if comment.Id == commentId {
					post.Comments = append(post.Comments[:j], post.Comments[j+1:]...)
					pp.Data[i] = post
					return nil
				}

			}

		}

	}
	return errors.New("this id doesnot exist")
}
