package auth

import (
	"errors"
	"net/http"
	"othnx/internal/domain"

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
		c1, c2, err := m.Service.CheckAndRefreshToken(accessToken.Value, refreshToken.Value)
		if c1 != nil && c2 != nil {
			c.SetCookie(c1)
			c.SetCookie(c2)
		}
		if err != nil {
			if err.Error() == domain.ErrInvalidAccessToken.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrInvalidAccessToken.Error()})
			}
			if err.Error() == domain.ErrRefreshTokenExpired.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrRefreshTokenExpired.Error()})
			}
			if err.Error() == domain.ErrUserDeleted.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrUserDeleted.Error()})
			}
			if err.Error() == domain.ErrUserBlocked.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrUserBlocked.Error()})
			}
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrGeneric.Error()})
		}
		return next(c)
	}
}

func (m AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return errors.New("not_implemented")
	}
}
