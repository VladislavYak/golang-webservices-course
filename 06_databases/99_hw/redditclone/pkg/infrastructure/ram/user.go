package ram

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/user"
)

var _ user.UserRepository = new(UserRepo)

type UserRepo struct {
	Users *[]user.User
	*sync.Mutex
	lastID int
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		&[]user.User{},
		&sync.Mutex{},
		0,
	}
}

func (ur *UserRepo) Create(ctx context.Context, User *user.User, Password string) (*user.User, error) {
	ur.Mutex.Lock()
	defer ur.Mutex.Unlock()

	if ur.isUserExists(User) {
		return nil, errors.New("this user already exists")
	}

	idStr := strconv.Itoa(ur.lastID)

	User = User.WithID(idStr)

	ur.lastID++

	*ur.Users = append(*ur.Users, *User)

	return User, nil

}

func (ur *UserRepo) isUserExists(user *user.User) bool {

	for _, userIter := range *ur.Users {
		if userIter.Username == user.Username {
			return true
		}
	}
	return false
}

func (ur *UserRepo) GetUser(ctx context.Context, user *user.User) (*user.User, error) {

	for _, userIter := range *ur.Users {
		if user.Username == userIter.Username {
			return &userIter, nil
		}
	}

	return nil, errors.New("user not found")
}

func (ur *UserRepo) GetUserPassword(ctx context.Context, user *user.User) (string, error) {
	return "", nil
}

func (r *UserRepo) AddJWT(context.Context, string, string, time.Time, time.Time) error {
	return nil
}

func (r *UserRepo) ValidateJWT(ctx context.Context, Token string, expiredAt time.Time) error {
	return nil
}
