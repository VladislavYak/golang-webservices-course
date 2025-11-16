package handlers

import (
	"context"

	"github.com/VladislavYak/redditclone/pkg/domain/auth"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func getUserCtx(c echo.Context) context.Context {
	us := c.Get("user").(*jwt.Token)
	claims := us.Claims.(*auth.JwtCustomClaims)
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "Username", claims.Login)
	ctx = context.WithValue(ctx, "UserID", claims.UserID)
	return ctx
}
