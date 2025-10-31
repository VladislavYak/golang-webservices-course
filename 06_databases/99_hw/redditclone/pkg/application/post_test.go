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
			s.PostMockRepo.EXPECT().Vote(s.ctx, myId, 1).Return(postAfterUpvote, nil),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(postAfterUpdateScore, nil),
		)

		returnedPost, err := s.postImpl.Create(s.ctx, myPost)

		s.Require().NoError(err)
		s.Assert().Equal(postAfterUpdateScore, returnedPost)
		s.Assert().Equal(1, returnedPost.Score)
	})

	t.Run("error on AddPost - no further calls", func(t *testing.T) {

		s.PostMockRepo.EXPECT().AddPost(s.ctx, myPost).Return(nil, mockErr).Times(1)
		s.PostMockRepo.EXPECT().Vote(gomock.Any(), gomock.Any(), 1).Times(0)
		s.PostMockRepo.EXPECT().UpdateScore(gomock.Any(), gomock.Any()).Times(0)

		returnedPost, err := s.postImpl.Create(s.ctx, myPost)

		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

	t.Run("error on Upvote - no UpdateScore call", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().AddPost(s.ctx, myPost).Return(postAfterAddPost, nil).Times(1),
			s.PostMockRepo.EXPECT().Vote(s.ctx, myId, 1).Return(nil, mockErr).Times(1),
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
			s.PostMockRepo.EXPECT().Vote(s.ctx, myId, 1).Return(postAfterUpvote, nil).Times(1),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(nil, mockErr).Times(1),
		)
		returnedPost, err := s.postImpl.Create(s.ctx, myPost)
		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

}

func (s *PostServiceTestSuite) TestGetAll() {
	t := s.T()

	tests := []struct {
		name    string
		setup   func()
		call    func() (interface{}, error)
		wantErr bool
		check   func(result interface{}, err error)
	}{
		{
			name: "GetAll - success",
			setup: func() {
				s.PostMockRepo.EXPECT().GetAllPosts(s.ctx).Return([]*post.Post{{Id: "1"}}, nil)
			},
			call: func() (interface{}, error) {
				return s.postImpl.GetAll(s.ctx)
			},
			wantErr: false,
		},
		{
			name: "GetAll - error",
			setup: func() {
				s.PostMockRepo.EXPECT().GetAllPosts(s.ctx).Return(nil, errors.New("db"))
			},
			call: func() (interface{}, error) {
				return s.postImpl.GetAll(s.ctx)
			},
			wantErr: true,
			check: func(_ interface{}, err error) {
				s.ErrorContains(err, "GetAll")
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result, err := tt.call()
			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
			if tt.check != nil {
				tt.check(result, err)
			}
		})
	}
}

func (s *PostServiceTestSuite) TestGetByID() {
	t := s.T()

	tests := []struct {
		name    string
		id      string
		setup   func(id string)
		call    func(id string) (interface{}, error)
		wantErr bool
		check   func(result interface{}, err error)
	}{
		{
			name: "GetByID - success",
			setup: func(id string) {
				s.PostMockRepo.EXPECT().GetPostByID(s.ctx, id).Return(&post.Post{Id: "123"}, nil)
			},
			call: func(id string) (interface{}, error) {
				return s.postImpl.GetByID(s.ctx, id)
			},
			wantErr: false,
		},
		{
			name: "GetByID - error",
			setup: func(id string) {
				s.PostMockRepo.EXPECT().GetPostByID(s.ctx, id).Return(nil, errors.New("db"))
			},
			call: func(id string) (interface{}, error) {
				return s.postImpl.GetByID(s.ctx, id)
			},
			wantErr: true,
			check: func(_ interface{}, err error) {
				s.ErrorContains(err, "GetByID")
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.id)
			result, err := tt.call(tt.id)
			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
			if tt.check != nil {
				tt.check(result, err)
			}
		})
	}

}

func (s *PostServiceTestSuite) TestGetPostsByCategoryName() {
	t := s.T()

	tests := []struct {
		name    string
		id      string
		setup   func(id string)
		call    func(id string) (interface{}, error)
		wantErr bool
		check   func(result interface{}, err error)
	}{
		{
			name: "GetPostsByCategoryName - success",
			setup: func(id string) {
				s.PostMockRepo.EXPECT().GetPostsByCategoryName(s.ctx, id).Return([]*post.Post{{Id: "123"}}, nil)
			},
			call: func(id string) (interface{}, error) {
				return s.postImpl.GetPostsByCategoryName(s.ctx, id)
			},
			wantErr: false,
		},
		{
			name: "GetPostsByCategoryName - error",
			setup: func(id string) {
				s.PostMockRepo.EXPECT().GetPostsByCategoryName(s.ctx, id).Return(nil, errors.New("db"))
			},
			call: func(id string) (interface{}, error) {
				return s.postImpl.GetPostsByCategoryName(s.ctx, id)
			},
			wantErr: true,
			check: func(_ interface{}, err error) {
				s.ErrorContains(err, "GetPostsByCategoryName")
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.id)
			result, err := tt.call(tt.id)
			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
			if tt.check != nil {
				tt.check(result, err)
			}
		})
	}

}

func (s *PostServiceTestSuite) TestGetPostsByUsername() {
	t := s.T()

	tests := []struct {
		name    string
		id      string
		setup   func(id string)
		call    func(id string) (interface{}, error)
		wantErr bool
		check   func(result interface{}, err error)
	}{
		{
			name: "GetPostsByUsername - success",
			setup: func(id string) {
				s.PostMockRepo.EXPECT().GetPostsByUsername(s.ctx, id).Return([]*post.Post{{Id: "123"}}, nil)
			},
			call: func(id string) (interface{}, error) {
				return s.postImpl.GetPostsByUsername(s.ctx, id)
			},
			wantErr: false,
		},
		{
			name: "GetPostsByUsername - error",
			setup: func(id string) {
				s.PostMockRepo.EXPECT().GetPostsByUsername(s.ctx, id).Return(nil, errors.New("db"))
			},
			call: func(id string) (interface{}, error) {
				return s.postImpl.GetPostsByUsername(s.ctx, id)
			},
			wantErr: true,
			check: func(_ interface{}, err error) {
				s.ErrorContains(err, "GetPostsByUsername")
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.id)
			result, err := tt.call(tt.id)
			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
			if tt.check != nil {
				tt.check(result, err)
			}
		})
	}

}

func (s *PostServiceTestSuite) TestUpvote() {
	t := s.T()

	myUsername := "testovTest"
	myId := "1123"
	myUser := user.NewUser(myUsername)
	myPost := postWithId(post.NewPost("music", "text", "", "boobar", "Boobarov", *myUser), myId)
	myVote := post.Vote{User: myUsername, VoteScore: 1}
	myScore := 1

	postAfterUpvote := postWithVote(myPost, &myVote)
	postAfterUpdateScore := postWithScore(postAfterUpvote, myScore)

	mockErr := errors.New("some error")

	t.Run("happy path", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().Vote(s.ctx, myId, 1).Return(postAfterUpvote, nil),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(postAfterUpdateScore, nil),
		)

		returnedPost, err := s.postImpl.Upvote(s.ctx, myId)

		s.Require().NoError(err)
		s.Assert().Equal(postAfterUpdateScore, returnedPost)
		s.Assert().Equal(1, returnedPost.Score)
	})

	t.Run("error on Vote - no further calls", func(t *testing.T) {

		s.PostMockRepo.EXPECT().Vote(gomock.Any(), gomock.Any(), 1).Return(nil, mockErr).Times(1)
		s.PostMockRepo.EXPECT().UpdateScore(gomock.Any(), gomock.Any()).Times(0)

		returnedPost, err := s.postImpl.Upvote(s.ctx, myId)

		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

	t.Run("error on UpdateScore", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().Vote(s.ctx, myId, 1).Return(postAfterUpvote, nil).Times(1),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(nil, mockErr).Times(1),
		)
		returnedPost, err := s.postImpl.Upvote(s.ctx, myId)
		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

}

func (s *PostServiceTestSuite) TestDownvote() {
	t := s.T()

	myUsername := "testovTest"
	myId := "1123"
	myUser := user.NewUser(myUsername)
	myPost := postWithId(post.NewPost("music", "text", "", "boobar", "Boobarov", *myUser), myId)
	myVote := post.Vote{User: myUsername, VoteScore: -1}
	myScore := -1

	postAfterDownvote := postWithVote(myPost, &myVote)
	postAfterUpdateScore := postWithScore(postAfterDownvote, myScore)

	mockErr := errors.New("some error")

	t.Run("happy path", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().Vote(s.ctx, myId, -1).Return(postAfterDownvote, nil),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(postAfterUpdateScore, nil),
		)

		returnedPost, err := s.postImpl.Downvote(s.ctx, myId)

		s.Require().NoError(err)
		s.Assert().Equal(postAfterUpdateScore, returnedPost)
		s.Assert().Equal(-1, returnedPost.Score)
	})

	t.Run("error on Vote - no further calls", func(t *testing.T) {

		s.PostMockRepo.EXPECT().Vote(gomock.Any(), gomock.Any(), -1).Return(nil, mockErr).Times(1)
		s.PostMockRepo.EXPECT().UpdateScore(gomock.Any(), gomock.Any()).Times(0)

		returnedPost, err := s.postImpl.Downvote(s.ctx, myId)

		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

	t.Run("error on UpdateScore", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().Vote(s.ctx, myId, -1).Return(postAfterDownvote, nil).Times(1),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(nil, mockErr).Times(1),
		)
		returnedPost, err := s.postImpl.Downvote(s.ctx, myId)
		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

}

func (s *PostServiceTestSuite) TestUnvote() {
	t := s.T()

	myUsername := "testovTest"
	myId := "1123"
	myUser := user.NewUser(myUsername)
	myPost := postWithId(post.NewPost("music", "text", "", "boobar", "Boobarov", *myUser), myId)
	myVote := post.Vote{User: myUsername, VoteScore: -1}
	myScore := -1

	postAfterDownvote := postWithVote(myPost, &myVote)
	postAfterUpdateScore := postWithScore(postAfterDownvote, myScore)

	mockErr := errors.New("some error")

	t.Run("happy path", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().Unvote(s.ctx, myId).Return(postAfterDownvote, nil),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(postAfterUpdateScore, nil),
		)

		returnedPost, err := s.postImpl.Unvote(s.ctx, myId)

		s.Require().NoError(err)
		s.Assert().Equal(postAfterUpdateScore, returnedPost)
		s.Assert().Equal(-1, returnedPost.Score)
	})

	t.Run("error on Vote - no further calls", func(t *testing.T) {

		s.PostMockRepo.EXPECT().Unvote(gomock.Any(), gomock.Any()).Return(nil, mockErr).Times(1)
		s.PostMockRepo.EXPECT().UpdateScore(gomock.Any(), gomock.Any()).Times(0)

		returnedPost, err := s.postImpl.Unvote(s.ctx, myId)

		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

	t.Run("error on UpdateScore", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().Unvote(s.ctx, myId).Return(postAfterDownvote, nil).Times(1),
			s.PostMockRepo.EXPECT().UpdateScore(s.ctx, myId).Return(nil, mockErr).Times(1),
		)
		returnedPost, err := s.postImpl.Unvote(s.ctx, myId)
		s.Require().Error(err)
		s.Assert().Nil(returnedPost)
		s.ErrorIs(err, mockErr)

	})

}

func (s *PostServiceTestSuite) TestDelete() {
	t := s.T()

	myUsername := "testovTest"
	myId := "1123"
	userID := "boobar"
	myUser := user.NewUser(myUsername)
	myUser.UserID = userID

	anotherUser := user.NewUser("another")
	anotherUser.UserID = userID
	// myPost := postWithId(post.NewPost("music", "text", "", "boobar", "Boobarov", *myUser), myId)

	toBeDeletedPost := post.NewPost("music", "text", "", "boobatbar", "youtube", *myUser)
	returnedPostGetPostById := post.NewPost("music", "text", "", "boobatbar", "youtube", *anotherUser)

	mockErr := errors.New("some error")
	t.Run("happy path", func(t *testing.T) {

		gomock.InOrder(
			s.PostMockRepo.EXPECT().GetPostByID(s.ctx, myId).Return(toBeDeletedPost, nil),
			s.PostMockRepo.EXPECT().DeletePost(s.ctx, myId).Return(toBeDeletedPost, nil),
		)

		_, err := s.postImpl.Delete(s.ctx, myId, myUser.UserID)

		s.Require().NoError(err)

	})

	t.Run("error on GetPostById - differentPostOwnerError", func(t *testing.T) {
		returnedPostGetPostById.Author.UserID = "invalid_user_id"

		s.PostMockRepo.EXPECT().GetPostByID(s.ctx, myId).Return(returnedPostGetPostById, nil).Times(1)
		s.PostMockRepo.EXPECT().DeletePost(s.ctx, myId).Times(0)

		pp, err := s.postImpl.Delete(s.ctx, myId, myUser.UserID)

		s.ErrorIs(err, post.DifferentPostOwnerError)
		s.Nil(pp)
	})

	t.Run("error on GetPostById", func(t *testing.T) {
		returnedPostGetPostById.Author.UserID = "invalid_user_id"

		s.PostMockRepo.EXPECT().GetPostByID(s.ctx, myId).Return(nil, mockErr)
		s.PostMockRepo.EXPECT().DeletePost(s.ctx, myId).Times(0)

		pp, err := s.postImpl.Delete(s.ctx, myId, myUser.UserID)

		s.ErrorIs(err, mockErr)
		s.Nil(pp)
	})

	t.Run("error on Delete", func(t *testing.T) {
		returnedPostGetPostById.Author.UserID = "invalid_user_id"

		s.PostMockRepo.EXPECT().GetPostByID(s.ctx, myId).Return(toBeDeletedPost, nil).Times(1)
		s.PostMockRepo.EXPECT().DeletePost(s.ctx, myId).Return(nil, mockErr).Times(1)

		pp, err := s.postImpl.Delete(s.ctx, myId, myUser.UserID)

		s.ErrorIs(err, mockErr)
		s.Nil(pp)
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

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}
