package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Login  string `json:"login"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(claims *JwtCustomClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
