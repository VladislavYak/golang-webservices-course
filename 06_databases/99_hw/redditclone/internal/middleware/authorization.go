package middleware

import (
	"errors"

	"github.com/VladislavYak/redditclone/pkg/application"
	"github.com/VladislavYak/redditclone/pkg/domain/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// yakovlev: prettify it, it looks awful right now
// lazy to prettify, it just works
func CustomAuth(config *echojwt.Config, authService *application.AuthImpl) echo.MiddlewareFunc {

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

			extractors, ceErr := echojwt.CreateExtractors(config.TokenLookup)

			if ceErr != nil {
				return ceErr
			}
			// var lastExtractorErr error
			// var lastTokenErr error

			for _, extractor := range extractors {
				auths, extrErr := extractor(c)
				if extrErr != nil {
					// lastExtractorErr = extrErr
					continue
				}
				for _, tokenString := range auths {

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

					if err != nil || !token.Valid {
						return err
					}

					err = authService.ValidateSession(c.Request().Context(), token.Raw, claims.ExpiresAt.Time)
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
