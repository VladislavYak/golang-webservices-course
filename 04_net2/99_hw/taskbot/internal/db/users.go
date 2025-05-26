package db

import (
	"errors"

	"github.com/VladislavYak/taskbot/models"
)

var USERS = NewUsers()

type Users struct {
	UsersDB []models.User
}

func NewUsers() *Users {
	return &Users{
		[]models.User{},
	}
}

func (u *Users) HasUser(user models.User) bool {
	for _, user_ := range u.UsersDB {
		if user_.UserID == user.UserID {
			return true
		}
	}
	return false
}

func (u *Users) AddUser(user models.User) error {
	if !u.HasUser(user) {
		u.UsersDB = append(u.UsersDB, user)
		return nil
	} else {
		return errors.New("this user user is already in db")
	}

}
