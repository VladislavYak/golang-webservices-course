package application

import (
	"context"
	"testing"

	"github.com/VladislavYak/redditclone/internal/mocks"
	"github.com/VladislavYak/redditclone/pkg/domain/comment"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/go-faster/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type CommentServiceTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	PostMockRepo    *mocks.MockPostRepository
	CommentMockRepo *mocks.MockCommentRepository
	commentImpl     *CommentImpl
	ctx             context.Context
}

func (s *CommentServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.CommentMockRepo = mocks.NewMockCommentRepository(s.ctrl)
	s.PostMockRepo = mocks.NewMockPostRepository(s.ctrl)

	s.commentImpl = &CommentImpl{PostRepo: s.PostMockRepo, CommentRepo: s.CommentMockRepo}
}

func (s *CommentServiceTestSuite) TearDownTest() {
	// gomock автоматически проверит все EXPECT
}

func (s *CommentServiceTestSuite) TestAddComment() {
	t := s.T()

	// Успешное добавление комментария
	t.Run("success - add comment and return updated post", func(t *testing.T) {
		postID := "post123"
		Co := &comment.Comment{
			Id:   "c1",
			Body: "mypost",
		}

		s.CommentMockRepo.EXPECT().
			AddComment(s.ctx, postID, Co).
			Return(nil)

		updatedPost := &post.Post{
			Id:       postID,
			Comments: []comment.Comment{*Co},
		}
		s.PostMockRepo.EXPECT().
			GetPostByID(s.ctx, postID).
			Return(updatedPost, nil)

		// ACT
		result, err := s.commentImpl.AddComment(s.ctx, postID, Co)

		// ASSERT
		s.NoError(err)
		s.Equal(updatedPost, result)
		s.Len(result.Comments, 1)
		s.Equal("mypost", result.Comments[0].Body)
	})

	t.Run("fail - AddComment post not found error", func(t *testing.T) {
		postID := "post123"
		comment := &comment.Comment{Id: "c1", Body: "Fail"}

		s.CommentMockRepo.EXPECT().
			AddComment(s.ctx, postID, comment).
			Return(post.PostNotFoundError)

		result, err := s.commentImpl.AddComment(s.ctx, postID, comment)

		s.ErrorContains(err, "post not found")
		s.Nil(result)
		// GetPostByID НЕ должен быть вызван
	})

	t.Run("fail - GetPostByID error after add", func(t *testing.T) {
		postID := "post123"
		comment := &comment.Comment{Id: "c1", Body: "Post not found"}

		s.CommentMockRepo.EXPECT().
			AddComment(s.ctx, postID, comment).
			Return(nil)

		s.PostMockRepo.EXPECT().
			GetPostByID(s.ctx, postID).
			Return(nil, errors.New("some error"))

		result, err := s.commentImpl.AddComment(s.ctx, postID, comment)

		s.ErrorContains(err, "some error")
		s.Nil(result)
	})
}

func (s *CommentServiceTestSuite) TestDeleteComment() {
	t := s.T()

	// Успешное удаление
	t.Run("success - delete comment and return updated post", func(t *testing.T) {
		postID := "post123"
		commentID := "c1"

		s.CommentMockRepo.EXPECT().
			DeleteComment(s.ctx, postID, commentID).
			Return(nil)

		updatedPost := &post.Post{
			Id:       postID,
			Comments: []comment.Comment{},
		}

		s.PostMockRepo.EXPECT().
			GetPostByID(s.ctx, postID).
			Return(updatedPost, nil)

		result, err := s.commentImpl.DeleteComment(s.ctx, postID, commentID)

		s.NoError(err)
		s.Equal(updatedPost, result)
		s.Empty(result.Comments)
	})

	// Ошибка при удалении
	t.Run("fail - DeleteComment post not found error", func(t *testing.T) {
		postID := "post123"
		commentID := "c1"

		s.CommentMockRepo.EXPECT().
			DeleteComment(s.ctx, postID, commentID).
			Return(comment.CommentNotFoundError)

		result, err := s.commentImpl.DeleteComment(s.ctx, postID, commentID)

		s.ErrorContains(err, "comment not found")
		s.Nil(result)
	})

	// Ошибка при получении поста после удаления
	t.Run("fail - GetPostByID error after delete", func(t *testing.T) {
		postID := "post123"
		commentID := "c1"

		s.CommentMockRepo.EXPECT().
			DeleteComment(s.ctx, postID, commentID).
			Return(nil)

		s.PostMockRepo.EXPECT().
			GetPostByID(s.ctx, postID).
			Return(nil, post.PostNotFoundError)

		result, err := s.commentImpl.DeleteComment(s.ctx, postID, commentID)

		s.ErrorContains(err, "post not found")
		s.Nil(result)
	})
}

func TestCommentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CommentServiceTestSuite))
}
