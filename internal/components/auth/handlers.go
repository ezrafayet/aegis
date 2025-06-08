package auth

import (
	"aegix/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandlersInterface interface {
	GetSession(c echo.Context) error
	Logout(c echo.Context) error
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
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrGeneric.Error()})
	} else {
		accessToken = cookie.Value
	}
	session, err := h.Service.GetSession(accessToken)
	if err != nil {
		if err.Error() == domain.ErrAccessTokenExpired.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrAccessTokenExpired.Error()})
		}
		if err.Error() == domain.ErrInvalidAccessToken.Error() {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": domain.ErrInvalidAccessToken.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": domain.ErrGeneric.Error()})
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
	c.SetCookie(&refreshCookie)
	c.SetCookie(&accessCookie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": domain.ErrGeneric.Error()})
	}
	return c.NoContent(http.StatusOK)
}
