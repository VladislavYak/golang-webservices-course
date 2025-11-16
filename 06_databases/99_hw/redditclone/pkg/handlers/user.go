package handlers

import (
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/domain/auth"
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

	if err := c.Bind(form); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	token, err := uh.Impl.Register(c.Request().Context(), form.Username, form.Password)
	if err != nil {

		if ve, ok := auth.AsValidationError(err); ok {
			return c.JSON(http.StatusUnprocessableEntity, ve)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
