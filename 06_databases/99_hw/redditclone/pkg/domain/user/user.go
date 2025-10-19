package user

import (
	"context"
	"time"
)

type User struct {
	Username string `json:"username"`
	// Password string
	UserID string `json:"id"`
}

func NewUser(Username string) *User {
	return &User{Username: Username}
}

func (u *User) WithID(Id string) *User {
	u.UserID = Id
	return u
}

// func (u *User) WithPassword(Password string) *User {
// 	u.Password = Password
// 	return u
// }

// func (u *User) GetPassword() string {
// 	return u.Password
// }

type UserRepository interface {
	Create(ctx context.Context, User *User, Password string) (*User, error)
	GetUser(ctx context.Context, User *User) (*User, error)
	GetUserPassword(ctx context.Context, user *User) (string, error)
	AddJWT(ctx context.Context, Token string, UserID string, IssuedAt time.Time, ExpiresAt time.Time) error
	ValidateJWT(ctx context.Context, Token string, ExpiresAt time.Time) error
}
