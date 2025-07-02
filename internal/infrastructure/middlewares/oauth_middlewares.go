package middlewares

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/pkg/apperrors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// some factory

type OAuthMiddlewaresInterface interface {
	CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc
}

type OAuthMiddlewares struct {
	Config  entities.Config
	Service primary.OAuthUseCasesExecutor
}

var _ OAuthMiddlewaresInterface = OAuthMiddlewares{}

func NewOAuthMiddlewares(c entities.Config) OAuthMiddlewares {
	return OAuthMiddlewares{
		Config: c,
	}
}

func (m OAuthMiddlewares) CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !m.Service.CheckAuthEnabled() {
			return c.JSON(http.StatusForbidden, map[string]string{"error": apperrors.ErrAuthMethodNotEnabled.Error()})
		}
		return next(c)
	}
}
