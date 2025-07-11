package middlewares

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/pkg/apperrors"
	"aegis/pkg/cookies"
	"errors"
	"net/http"

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

var _ AuthMiddlewareInterface = (*AuthMiddleware)(nil)

func NewAuthMiddleware(c entities.Config, s primary.UseCasesInterface) *AuthMiddleware {
	return &AuthMiddleware{
		Config:  c,
		Service: s,
	}
}

func (m AuthMiddleware) CheckAndRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessTokenValue := ""
		if accessToken, err := c.Cookie("access_token"); err == nil {
			accessTokenValue = accessToken.Value
		}
		refreshTokenValue := ""
		if refreshToken, err := c.Cookie("refresh_token"); err == nil {
			refreshTokenValue = refreshToken.Value
		}
		tokensPair, err := m.Service.CheckAndRefreshToken(accessTokenValue, refreshTokenValue, false)
		if tokensPair != nil {
			accessCookie := cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), m.Config)
			refreshCookie := cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), m.Config)
			c.SetCookie(&accessCookie)
			c.SetCookie(&refreshCookie)
		}
		if err != nil {
			if errors.Is(err, apperrors.ErrAccessTokenInvalid) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenInvalid.Error()})
			}
			if errors.Is(err, apperrors.ErrRefreshTokenExpired) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrRefreshTokenExpired.Error()})
			}
			if errors.Is(err, apperrors.ErrUserDeleted) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserDeleted.Error()})
			}
			if errors.Is(err, apperrors.ErrUserBlocked) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserBlocked.Error()})
			}
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
		}
		return next(c)
	}
}

func (m AuthMiddleware) CheckAndForceRefreshToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessTokenValue := ""
		if accessToken, err := c.Cookie("access_token"); err == nil {
			accessTokenValue = accessToken.Value
		}
		refreshTokenValue := ""
		if refreshToken, err := c.Cookie("refresh_token"); err == nil {
			refreshTokenValue = refreshToken.Value
		}
		tokensPair, err := m.Service.CheckAndRefreshToken(accessTokenValue, refreshTokenValue, true)
		if tokensPair != nil {
			accessCookie := cookies.NewAccessCookie(tokensPair.AccessToken, tokensPair.AccessTokenExpiresAt.Unix(), m.Config)
			refreshCookie := cookies.NewRefreshCookie(tokensPair.RefreshToken, tokensPair.RefreshTokenExpiresAt.Unix(), m.Config)
			c.SetCookie(&accessCookie)
			c.SetCookie(&refreshCookie)
		}
		if err != nil {
			if errors.Is(err, apperrors.ErrAccessTokenInvalid) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenInvalid.Error()})
			}
			if errors.Is(err, apperrors.ErrRefreshTokenExpired) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrRefreshTokenExpired.Error()})
			}
			if errors.Is(err, apperrors.ErrUserDeleted) {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrUserDeleted.Error()})
			}
			if errors.Is(err, apperrors.ErrUserBlocked) {
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
