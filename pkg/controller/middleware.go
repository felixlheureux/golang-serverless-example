package controller

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manta-coder/golang-serverless-example/pkg/auth"
	"github.com/manta-coder/golang-serverless-example/pkg/httperror"
)

func NewAuthenticator(secret string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.Claims{},
		SigningKey: []byte(secret),
		ErrorHandler: func(err error) error {
			return httperror.CoreUnauthorized(err)
		},
	})
}

func NewAuthMiddleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := getClaims(c)
			if c.QueryParam("ethereum_address") != claims.EthereumAddressHex {
				return httperror.CoreUnauthorized(fmt.Errorf("invalid address"))
			}
			return h(c)
		}
	}
}

func getClaims(c echo.Context) *auth.Claims {
	user := c.Get("user").(*jwt.Token)
	return user.Claims.(*auth.Claims)
}
