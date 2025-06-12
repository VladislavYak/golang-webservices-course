package handlers

import (
	"fmt"
	"net/http"

	"github.com/VladislavYak/redditclone/pkg/user"
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

	user, err := lh.UserRepo.GetUser(&user.User{Username: form.Username, Password: form.Password})
	if err != nil {
		echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if user.Password != form.Password {
		echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	// if err := rh.UserRepo.AddUser(&user.User{Username: form.Username, Password: form.Password}); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, err)
	// }

	fmt.Println("f", form)

	// c.String(http.StatusCreated, "test")

	// fmt.Println(rh.UserRepo.Users)
	return nil
}
