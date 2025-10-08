package ram

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/VladislavYak/redditclone/pkg/domain/user"
)

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

func (ur *UserRepo) AddUser(user *user.User) (string, error) {
	ur.Mutex.Lock()
	defer ur.Mutex.Unlock()

	fmt.Println("ur.Users AddUser", ur.Users)
	if ur.isUserExists(user) {
		return "", errors.New("this user already exists")
	}

	idStr := strconv.Itoa(ur.lastID)

	user = user.WithID(idStr)

	ur.lastID++

	*ur.Users = append(*ur.Users, *user)

	fmt.Println("users", ur.Users)
	return idStr, nil

}

func (ur *UserRepo) isUserExists(user *user.User) bool {

	fmt.Println("ur.Users")
	fmt.Printf("isUserExists ur.Users %+v\n", ur.Users)
	for _, userIter := range *ur.Users {
		if userIter.Username == user.Username {
			return true
		}
	}
	return false
}

func (ur *UserRepo) GetUser(user *user.User) (*user.User, error) {
	fmt.Println("inside GetUser, before loop")
	fmt.Printf("in use %+v\n", user)
	fmt.Printf("ur.Users %+v\n", ur.Users)
	for _, userIter := range *ur.Users {
		fmt.Println("userIter", userIter)
		fmt.Printf("userIter %+v\n", userIter)
		if user.Username == userIter.Username {
			return &userIter, nil
		}
	}

	fmt.Println("inside GetUser after loop")

	return nil, errors.New("user not found")
}
