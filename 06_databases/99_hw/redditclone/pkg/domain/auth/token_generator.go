package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(Claims *JwtCustomClaims, Secret string) (string, error) {
	// yakovlev: temp hardcoding

	Secret = "secret"

	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return Token.SignedString([]byte(Secret))
}
