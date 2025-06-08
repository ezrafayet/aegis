package auth

import (
	"aegix/internal/domain"
	"net/http"

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
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrGeneric.Error()})
		}
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrGeneric.Error()})
		}
		c1, c2, setC, err := m.Service.CheckAndRefreshToken(accessToken.Value, refreshToken.Value)
		if err != nil {
			if err.Error() == domain.ErrInvalidAccessToken.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrInvalidAccessToken.Error()})
			}
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrGeneric.Error()})
		}
		if setC {
			c.SetCookie(&c1)
			c.SetCookie(&c2)
		}
		return next(c)
	}
}

func (m AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
