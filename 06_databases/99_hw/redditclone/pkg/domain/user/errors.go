package user

import "github.com/go-faster/errors"

var (
	UserNotExistsError     = errors.New("user not found")
	UserAlreadyExistsError = errors.New("Username already exists")
)
