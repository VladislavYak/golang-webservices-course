package handlers

import (
	"net/http"
	"time"

	"github.com/VladislavYak/redditclone/pkg/user"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterHandler struct {
	UserRepo user.UserRepo
}

func (rh *RegisterHandler) Register(c echo.Context) error {

	form := &RegisterForm{}

	if err := c.Bind(form); err != nil {
		return err
	}

	if err := rh.UserRepo.AddUser(&user.User{Username: form.Username, Password: form.Password}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Set custom claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	}

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
