package handlers

import (
	"fmt"
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/labstack/echo/v4"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	Impl application.UserInterface
}

func (uh *UserHandler) Login(c echo.Context) error {
	form := &LoginForm{}

	if err := c.Bind(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	token, err := uh.Impl.Login(c.Request().Context(), form.Username, form.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (uh *UserHandler) Register(c echo.Context) error {

	form := &RegisterForm{}

	fmt.Println("before logging")

	if err := c.Bind(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	fmt.Println("form.Username", form.Username)
	token, err := uh.Impl.Register(c.Request().Context(), form.Username, form.Password)
	if err != nil {
		fmt.Println("i was here")
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	fmt.Println("token Register", token)

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
