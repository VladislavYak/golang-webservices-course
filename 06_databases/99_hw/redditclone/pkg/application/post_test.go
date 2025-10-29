package application

import (
	"context"
	"testing"

	"github.com/VladislavYak/redditclone/mocks"
	"github.com/VladislavYak/redditclone/pkg/domain/post"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"github.com/go-faster/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PostServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	PostMockRepo *mocks.MockPostRepository
	postImpl     *PostImpl
	ctx          context.Context
}

func (s *PostServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.PostMockRepo = mocks.NewMockPostRepository(s.ctrl)

	s.postImpl = &PostImpl{repo: s.PostMockRepo}
}

func (s *PostServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *PostServiceTestSuite) TestCreate() {
	t := s.T()

	myUsername := "testovTest"
	myId := "1123"
	myUser := user.NewUser(myUsername)
	myPost := post.NewPost("music", "text", "", "boobar", "Boobarov", *myUser)
	myVote := post.Vote{User: myUsername, VoteScore: 1}
	myScore := 1

	postAfterAddPost := postWithId(myPost, myId)
	postAfterUpvote := postWithVote(postAfterAddPost, &myVote)
	postAfterUpdateScore := postWithScore(postAfterUpvote, myScore)

	mockErr := errors.New("some error")

	t.Run("happy path", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().AddPost(s.ctx, myPost).Return(postAfterAddPost, nil),
			s.PostMockRepo.EXPECT().Upvote(s.ctx, myId).Return(postAfterUpvote, nil),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(postAfterUpdateScore, nil),
		)

		returnedPost, err := s.postImpl.Create(s.ctx, myPost)

		s.Require().NoError(err)
		s.Assert().Equal(postAfterUpdateScore, returnedPost)
		s.Assert().Equal(1, returnedPost.Score)
	})

	t.Run("error on AddPost - no further calls", func(t *testing.T) {

		s.PostMockRepo.EXPECT().AddPost(s.ctx, myPost).Return(nil, mockErr).Times(1)
		s.PostMockRepo.EXPECT().Upvote(gomock.Any(), gomock.Any()).Times(0)
		s.PostMockRepo.EXPECT().UpdateScore(gomock.Any(), gomock.Any()).Times(0)

		returnedPost, err := s.postImpl.Create(s.ctx, myPost)

		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

	t.Run("error on Upvote - no UpdateScore call", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().AddPost(s.ctx, myPost).Return(postAfterAddPost, nil).Times(1),
			s.PostMockRepo.EXPECT().Upvote(s.ctx, myId).Return(nil, mockErr).Times(1),
		)

		s.PostMockRepo.EXPECT().UpdateScore(gomock.Any(), gomock.Any()).Times(0) // Не вызван

		returnedPost, err := s.postImpl.Create(s.ctx, myPost)
		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)
	})

	t.Run("error on UpdateScore", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().AddPost(s.ctx, myPost).Return(postAfterAddPost, nil).Times(1),
			s.PostMockRepo.EXPECT().Upvote(s.ctx, myId).Return(postAfterUpvote, nil).Times(1),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(nil, mockErr).Times(1),
		)
		returnedPost, err := s.postImpl.Create(s.ctx, myPost)
		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

}

func postWithId(pst *post.Post, id string) *post.Post {
	clone := *pst
	clone.Id = id
	return &clone
}

func postWithVote(pst *post.Post, vote *post.Vote) *post.Post {
	clone := *pst
	// yakovlev: be careful here
	clone.Votes = append([]post.Vote(nil), *vote)
	return &clone
}

func postWithScore(pst *post.Post, score int) *post.Post {

	clone := *pst

	clone.Score = score

	return &clone
}
