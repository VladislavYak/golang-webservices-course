package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

// yakovlev: вынести из domain куда-то в другое место
func GenerateJWTToken(Claims *JwtCustomClaims, Secret string) (string, error) {

	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)

	SignedToken, err := Token.SignedString([]byte(Secret))
	if err != nil {
		return "", InvalidTokenError
	}
	return SignedToken, nil
}
