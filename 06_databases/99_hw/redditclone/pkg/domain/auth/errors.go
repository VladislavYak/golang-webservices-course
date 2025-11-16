package auth

import (
	"encoding/json"
	"errors"
)

var (
	ExpiredTokenError    = errors.New("token expired")
	InvalidTokenError    = errors.New("invalid token")
	InvalidPasswordError = errors.New("invalid password")
)

type ValidationError struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value"`
	Msg      string `json:"msg"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve ValidationErrors) Error() string {
	b, _ := json.Marshal(ve)
	return string(b)
}

// Optional: Helper to check if an error is ValidationErrors
func IsValidationError(err error) bool {
	_, ok := err.(ValidationErrors)
	return ok
}

func AsValidationError(err error) (ValidationErrors, bool) {
	ve, ok := err.(ValidationErrors)
	return ve, ok
}
