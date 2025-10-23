package application

import (
	"context"
	"testing"

	"github.com/VladislavYak/redditclone/mocks"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"github.com/go-faster/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	UsermockRepo *mocks.MockUserRepository
	AuthmockRepo *mocks.MockAuthRepository
	userImpl     *UserImpl
	ctx          context.Context
}

func (s *UserServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.UsermockRepo = mocks.NewMockUserRepository(s.ctrl)
	s.AuthmockRepo = mocks.NewMockAuthRepository(s.ctrl)
	s.userImpl = &UserImpl{
		ur:        s.UsermockRepo,
		ar:        s.AuthmockRepo,
		JWTSecret: "secret",
	}
}

func (s *UserServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *UserServiceTestSuite) TestRegister() {
	t := s.T()

	t.Run("user already exists error", func(t *testing.T) {

		User := user.NewUser("Vlad")

		s.UsermockRepo.EXPECT().GetUser(s.ctx, User).Return(&user.User{Username: "Vlad", UserID: "boobar"}, nil)

		_, err := s.userImpl.Register(s.ctx, "Vlad", "pass")

		s.Assert().ErrorIs(err, user.UserAlreadyExistsError)

	})

	t.Run("no error", func(t *testing.T) {

		User := user.NewUser("Vlad")

		s.UsermockRepo.EXPECT().
			GetUser(s.ctx, User).
			Return(nil, user.UserNotExistsError)
		s.UsermockRepo.EXPECT().
			Create(s.ctx, User, "pass").
			Return(&user.User{Username: "Vlad", UserID: "1234"}, nil)
		s.AuthmockRepo.EXPECT().
			AddJWT(s.ctx, gomock.Any(), "1234", gomock.Any(), gomock.Any()).
			Return(nil)

		token, err := s.userImpl.Register(s.ctx, "Vlad", "pass")

		s.NoError(err)
		s.NotEmpty(token)

	})

	t.Run("db error", func(t *testing.T) {
		expected := ""

		User := user.NewUser("Vlad")

		s.UsermockRepo.EXPECT().
			GetUser(s.ctx, User).
			Return(nil, user.UserNotExistsError)
		s.UsermockRepo.EXPECT().
			Create(s.ctx, User, "pass").
			Return(nil, errors.New("err"))

		token, err := s.userImpl.Register(s.ctx, "Vlad", "pass")

		s.Error(err)
		s.Equal(token, expected)

	})

}

func TestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
