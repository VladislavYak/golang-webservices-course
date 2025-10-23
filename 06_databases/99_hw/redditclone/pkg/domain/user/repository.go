package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, User *User, Password string) (*User, error)
	GetUser(ctx context.Context, User *User) (*User, error)
	GetUserPassword(ctx context.Context, user *User) (string, error)
}
