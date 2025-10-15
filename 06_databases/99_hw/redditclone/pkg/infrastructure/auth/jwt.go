package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Login  string `json:"login"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(Claims *JwtCustomClaims, Secret string) (string, error) {
	// yakovlev: temp hardcoding

	Secret = "secret"

	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return Token.SignedString([]byte(Secret))
}
