package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"othnx/internal/domain"
	"othnx/pkg/apperrors"
)

type AuthHandlersInterface interface {
	GetSession(c echo.Context) error
	Logout(c echo.Context) error
	DoNothing(c echo.Context) error
}

type AuthHandlers struct {
	Config  domain.Config
	Service AuthServiceInterface
}

var _ AuthHandlersInterface = &AuthHandlers{}

func NewAuthHandlers(c domain.Config, s AuthServiceInterface) AuthHandlers {
	return AuthHandlers{
		Config:  c,
		Service: s,
	}
}

func (h AuthHandlers) GetSession(c echo.Context) error {
	var accessToken string
	if cookie, err := c.Cookie("access_token"); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrGeneric.Error()})
	} else {
		accessToken = cookie.Value
	}
	session, err := h.Service.GetSession(accessToken)
	if err != nil {
		if err.Error() == apperrors.ErrAccessTokenExpired.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenExpired.Error()})
		}
		if err.Error() == apperrors.ErrAccessTokenInvalid.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": apperrors.ErrAccessTokenInvalid.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	return c.JSON(http.StatusOK, session)
}

func (h AuthHandlers) Logout(c echo.Context) error {
	var refreshToken string
	if cookie, err := c.Cookie("refresh_token"); err != nil {
		refreshToken = ""
	} else {
		refreshToken = cookie.Value
	}
	accessCookie, refreshCookie, err := h.Service.Logout(refreshToken)
	if accessCookie != nil && refreshCookie != nil {
		c.SetCookie(accessCookie)
		c.SetCookie(refreshCookie)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": apperrors.ErrGeneric.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h AuthHandlers) DoNothing(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
