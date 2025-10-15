package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/user"
	"github.com/VladislavYak/redditclone/pkg/infrastructure/auth"
	jwt "github.com/golang-jwt/jwt/v5"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	jwt.RegisteredClaims
}

type UserInterface interface {
	Login(ctx context.Context, Login, Password string) (string, error)
	Register(ctx context.Context, Login, Password string) (string, error)
}

var _ UserInterface = new(UserImpl)

type UserImpl struct {
	ur        user.UserRepository
	JWTSecret string
}

func NewUserImpl(repo user.UserRepository) *UserImpl {
	// yakovlev: JWTSecret which?
	return &UserImpl{ur: repo}
}

func (ui *UserImpl) Register(ctx context.Context, Login, Password string) (string, error) {

	// yakovlev: это нужно как-то добавить
	// if _, err := ui.ur.GetUser(ctx, User); err == nil {
	//     return "", echo.NewHTTPError(http.StatusBadRequest, "user with this login already exists")
	// }

	u := user.NewUser(Login)

	fmt.Println("before ui.ur.Cerate")
	u, err := ui.ur.Create(ctx, u, Password)
	if err != nil {
		return "", err
	}

	// Генерируем JWT
	claims := &auth.JwtCustomClaims{
		Login:  Login,
		UserID: u.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token, err := auth.GenerateJWTToken(claims, ui.JWTSecret)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (ui *UserImpl) Login(ctx context.Context, Login, Password string) (string, error) {

	User, err := ui.ur.GetUser(ctx, user.NewUser(Login))
	if err != nil {
		return "", err
	}

	actualPass, err := ui.ur.GetUserPassword(ctx, User)
	if err != nil {
		return "", err
	}

	if actualPass != Password {
		return "", errors.New("invalid password")
	}

	claims := &auth.JwtCustomClaims{
		Login:  Login,
		UserID: User.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token, err := auth.GenerateJWTToken(claims, ui.JWTSecret)

	return token, nil

}
