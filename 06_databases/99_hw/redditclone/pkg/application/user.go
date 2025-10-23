package application

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"

	"github.com/VladislavYak/redditclone/pkg/domain/auth"
	"github.com/VladislavYak/redditclone/pkg/domain/user"
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

type UserInterface interface {
	Login(ctx context.Context, Login, Password string) (string, error)
	Register(ctx context.Context, Login, Password string) (string, error)
	// ValidateSession(ctx context.Context, Token string, ExpiresAt time.Time) error
}

var _ UserInterface = new(UserImpl)

type UserImpl struct {
	ur user.UserRepository
	ar auth.AuthRepository
	// lazy now
	JWTSecret string
}

func NewUserImpl(repo user.UserRepository, AuthRepo auth.AuthRepository) *UserImpl {
	// yakovlev: JWTSecret which?
	return &UserImpl{ur: repo, ar: AuthRepo}
}

func (ui *UserImpl) Register(ctx context.Context, Login, Password string) (string, error) {
	const op = "Register"

	u := user.NewUser(Login)

	if _, err := ui.ur.GetUser(ctx, u); !errors.Is(err, user.UserNotExistsError) {
		return "", user.UserAlreadyExistsError
	}

	fmt.Println("before ui.ur.Cerate")
	u, err := ui.ur.Create(ctx, u, Password)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(15 * time.Minute) // Shortened lifetime for security

	// Генерируем JWT
	Claims := &auth.JwtCustomClaims{
		Login:  Login,
		UserID: u.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	Token, err := auth.GenerateJWTToken(Claims, ui.JWTSecret)

	if err != nil {
		return "", errors.Wrap(err, op)
	}

	if err = ui.ar.AddJWT(ctx, Token, Claims.UserID, Claims.IssuedAt.Time, Claims.ExpiresAt.Time); err != nil {
		return "", errors.Wrap(err, op)
	}

	return Token, nil
}

func (ui *UserImpl) Login(ctx context.Context, Login, Password string) (string, error) {
	const op = "Login"

	User, err := ui.ur.GetUser(ctx, user.NewUser(Login))
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	actualPass, err := ui.ur.GetUserPassword(ctx, User)
	if err != nil {
		return "", errors.Wrap(err, op)
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

	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return token, nil

}

// // need to be moved to auth
// func (ui *UserImpl) ValidateSession(ctx context.Context, Token string, ExpiresAt time.Time) error {
// 	const op = "ValidateSession"
// 	err := ui.ur.ValidateJWT(ctx, Token, ExpiresAt)

// 	if err != nil {
// 		return errors.Wrap(err, op)
// 	}

// 	return nil
// }
