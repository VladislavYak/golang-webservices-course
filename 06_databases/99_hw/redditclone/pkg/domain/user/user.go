package user

import "context"

type User struct {
	Username string `json:"username"`
	Password string
	UserID   string `json:"id"`
}

func NewUser(Username string) *User {
	return &User{Username: Username}
}

func (u *User) WithID(Id string) *User {
	u.UserID = Id
	return u
}

func (u *User) WithPassword(Password string) *User {
	u.Password = Password
	return u
}

func (u *User) GetPassword() string {
	return u.Password
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, user *User) (*User, error)
}
