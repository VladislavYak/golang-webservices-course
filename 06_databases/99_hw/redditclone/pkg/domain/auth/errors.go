package auth

import "errors"

var (
	ExpiredTokenError    = errors.New("token expired")
	InvalidTokenError    = errors.New("invalid token")
	InvalidPasswordError = errors.New("invalid password")
)
