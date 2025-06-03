package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RegisterForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	form := &RegisterForm{}

	if err := c.Bind(form); err != nil {
		return err
	}

	fmt.Println("f", form)

	c.String(http.StatusCreated, "test")
	return nil
}
