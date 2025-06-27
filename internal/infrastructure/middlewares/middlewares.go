package middlewares

import (
	"errors"
	"net/http"
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/primary"
	"othnx/pkg/apperrors"
	"othnx/pkg/cookies"

	"github.com/labstack/echo/v4"
)

type AuthMiddlewareInterface interface {
	CheckToken(next echo.HandlerFunc) echo.HandlerFunc
	CheckAndRefreshToken(next echo.HandlerFunc) echo.HandlerFunc
	CheckAndForceRefreshToken(next echo.HandlerFunc) echo.HandlerFunc
}

type AuthMiddleware struct {
	Config  entities.Config
	Service primary.UseCasesInterface
}

var _ AuthMiddlewareInterface = &AuthMiddleware{}

func NewAuthMiddleware(c entities.Config, s primary.UseCasesInterface) AuthMiddleware {
	return AuthMiddleware{
		Config:  c,
		Service: s,
	}
}

func (m AuthMiddleware) CheckAndRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
		}
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
		}
		tokensPair, err := m.Service.CheckAndRefreshToken(accessToken.Value, refreshToken.Value, false)
		if tokensPair != nil {
			cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), m.Config)
			cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), m.Config)
		}
		if err != nil {
			if err.Error() == apperrors.ErrAccessTokenInvalid.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenInvalid.Error()})
			}
			if err.Error() == apperrors.ErrRefreshTokenExpired.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrRefreshTokenExpired.Error()})
			}
			if err.Error() == apperrors.ErrUserDeleted.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserDeleted.Error()})
			}
			if err.Error() == apperrors.ErrUserBlocked.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserBlocked.Error()})
			}
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
		}
		return next(c)
	}
}

func (m AuthMiddleware) CheckAndForceRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
		}
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
		}
		tokensPair, err := m.Service.CheckAndRefreshToken(accessToken.Value, refreshToken.Value, true)
		if tokensPair != nil {
			cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), m.Config)
			cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), m.Config)
		}
		if err != nil {
			if err.Error() == apperrors.ErrAccessTokenInvalid.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenInvalid.Error()})
			}
			if err.Error() == apperrors.ErrRefreshTokenExpired.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrRefreshTokenExpired.Error()})
			}
			if err.Error() == apperrors.ErrUserDeleted.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserDeleted.Error()})
			}
			if err.Error() == apperrors.ErrUserBlocked.Error() {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserBlocked.Error()})
			}
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
		}
		return next(c)
	}
}

func (m AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return errors.New("not_implemented")
	}
}
