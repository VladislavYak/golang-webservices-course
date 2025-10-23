package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Login  string `json:"login"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
