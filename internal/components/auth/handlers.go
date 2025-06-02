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
	return c.NoContent(http.StatusOK)
}

func (h AuthHandlers) Logout(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
