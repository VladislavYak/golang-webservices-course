package middleware

import (
	"errors"
	"fmt"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/infrastructure/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// yakovlev: prettify it, it looks awful right now
func CustomAuth(config *echojwt.Config, userService application.UserInterface) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if config.Skipper == nil {
				config.Skipper = middleware.DefaultSkipper
			}
			if config.ContextKey == "" {
				config.ContextKey = "user"
			}
			if config.TokenLookup == "" && len(config.TokenLookupFuncs) == 0 {
				config.TokenLookup = "header:Authorization:Bearer "
			}
			if config.SigningMethod == "" {
				config.SigningMethod = echojwt.AlgorithmHS256
			}

			if config.NewClaimsFunc == nil {
				config.NewClaimsFunc = func(c echo.Context) jwt.Claims {
					return jwt.MapClaims{}
				}
			}
			if config.SigningKey == nil && len(config.SigningKeys) == 0 && config.KeyFunc == nil && config.ParseTokenFunc == nil {
				return errors.New("jwt middleware requires signing key")
			}
			// if config.KeyFunc == nil {
			// 	config.KeyFunc = config.defaultKeyFunc
			// }
			// if config.ParseTokenFunc == nil {
			// 	config.ParseTokenFunc = echojwt.defaultParseTokenFunc
			// }

			// if len(config.TokenLookupFuncs) > 0 {
			// 	extractors = append(config.TokenLookupFuncs, extractors...)
			// }

			fmt.Println("im inside custom auth")

			fmt.Println("config.TokenLookup", config.TokenLookup)

			extractors, ceErr := echojwt.CreateExtractors(config.TokenLookup)

			fmt.Println("extracted extractor")

			if ceErr != nil {
				fmt.Println("ceErr CreateExtractors", ceErr)
				return ceErr
			}
			// var lastExtractorErr error
			// var lastTokenErr error

			fmt.Println("len(extractors)", len(extractors))

			for _, extractor := range extractors {
				auths, extrErr := extractor(c)
				fmt.Println("len(auths)", len(auths))
				fmt.Println("extrErr", extrErr)
				if extrErr != nil {
					// lastExtractorErr = extrErr
					continue
				}
				for _, tokenString := range auths {
					fmt.Println("tokenString", tokenString)

					hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
						// method, ok := token.Method.(*jwt.SigningMethodHMAC)
						// if !ok || method.Alg() != "HS256" {
						// 	return nil, fmt.Errorf("bad sign method")
						// }
						return []byte("secret"), nil
					}

					// token, err := jwt.Parse(tokenString, hashSecretGetter)

					claims := auth.JwtCustomClaims{}
					token, err := jwt.ParseWithClaims(tokenString, &claims, hashSecretGetter)

					fmt.Println("token (after Parse)", token)
					fmt.Println("err (after Parse)", err)
					if err != nil || !token.Valid {
						return err
					}

					fmt.Println("before running ValidateSession")
					err = userService.ValidateSession(c.Request().Context(), token.Raw, claims.ExpiresAt.Time)
					if err != nil {
						// yakovlev: пока что хз как тут ошибки обарабывать, errors.Wrap или ХТТПШные?
						return err
					}

					// Store user information from token into context.
					c.Set(config.ContextKey, token)
					if config.SuccessHandler != nil {
						config.SuccessHandler(c)
					}
					return next(c)
				}
			}

			return nil

		}
	}
}
