package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/VladislavYak/redditclone/pkg/domain/user"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginHandler struct {
	UserRepo user.UserRepo
}

func (lh *LoginHandler) Login(c echo.Context) error {
	form := &LoginForm{}

	if err := c.Bind(form); err != nil {
		return err
	}
	user, err := lh.UserRepo.GetUser(user.NewUser(form.Username).WithPassword(form.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if user.GetPassword() != form.Password {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
	}

	ctx := c.Request().Context()
	context.WithValue(ctx, "user", form.Username)
	c.Request().WithContext(ctx)

	claims := &JwtCustomClaims{
		form.Username,
		user.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// YAKOVLEV: NOT TESTED. IS TOKEN VALID? PASSWORD?

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
