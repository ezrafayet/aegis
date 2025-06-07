package auth

import (
	"aegix/internal/domain"

	"github.com/labstack/echo/v4"
)

type AuthMiddlewareInterface interface {
	CheckToken(next echo.HandlerFunc) echo.HandlerFunc
	CheckAndRefreshToken(next echo.HandlerFunc) echo.HandlerFunc
}

type AuthMiddleware struct {
	Config  domain.Config
	Service AuthServiceInterface
}

var _ AuthMiddlewareInterface = &AuthMiddleware{}

func NewAuthMiddleware(c domain.Config, s AuthServiceInterface) AuthMiddleware {
	return AuthMiddleware{
		Config:  c,
		Service: s,
	}
}

func (m AuthMiddleware) CheckAndRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (m AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
