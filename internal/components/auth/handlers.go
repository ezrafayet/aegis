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
	Config domain.Config
	Service AuthServiceInterface
}

var _ AuthHandlersInterface = &AuthHandlers{}

func NewAuthHandlers(c domain.Config, s AuthServiceInterface) AuthHandlers {
	return AuthHandlers{
		Config: c,
		Service: s,
	}
}

func (h AuthHandlers) GetSession(c echo.Context) error {
	cookie, err := c.Cookie("access_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "an error occured"})
	}
	accessToken := cookie.Value
	session, err := h.Service.GetSession(accessToken)
	if err != nil {
		if err.Error() == "access_token_expired" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "access_token_expired"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an error occured"})
	}
	return c.JSON(http.StatusOK, session)
}

func (h AuthHandlers) Logout(c echo.Context) error {
	var refreshToken string
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		refreshToken = ""
	} else {
		refreshToken = cookie.Value
	}
	accessCookie, refreshCookie, err := h.Service.Logout(refreshToken)
	if err != nil {
		c.SetCookie(&refreshCookie)
		c.SetCookie(&accessCookie)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "an error occured"})
	}
	c.SetCookie(&refreshCookie)
	c.SetCookie(&accessCookie)
	return c.NoContent(http.StatusOK)
}
