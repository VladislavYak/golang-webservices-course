package user

import (
	"errors"
	"strconv"
	"sync"
)

type UserRepo struct {
	Users []*User
	*sync.Mutex
	lastID int
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		[]*User{},
		&sync.Mutex{},
		0,
	}
}

func (ur *UserRepo) AddUser(user *User) error {
	ur.Mutex.Lock()
	defer ur.Mutex.Unlock()

	if ur.isUserExists(user) {
		return errors.New("this user already exists")
	}

	idStr := strconv.Itoa(ur.lastID)
	user.UserID = idStr
	ur.lastID++

	ur.Users = append(ur.Users, user)
	return nil

}

func (ur *UserRepo) isUserExists(user *User) bool {
	for _, userIter := range ur.Users {
		if userIter.Username == user.Username {
			return true
		}
	}
	return false
}

func (ur *UserRepo) GetUser(user *User) (*User, error) {
	for _, userIter := range ur.Users {
		if user.Username == userIter.Username {
			return userIter, nil
		}
	}

	return nil, errors.New("user not found")
}
