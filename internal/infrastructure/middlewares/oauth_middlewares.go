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

type OAuthGithubMiddlewares struct {
	Config  entities.Config
	Service primary.OAuthUseCasesExecutor
}

var _ OAuthMiddlewaresInterface = OAuthGithubMiddlewares{}

func NewOAuthGithubMiddlewares(c entities.Config) OAuthGithubMiddlewares {
	return OAuthGithubMiddlewares{
		Config: c,
	}
}

func (m OAuthGithubMiddlewares) CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !m.Service.CheckAuthEnabled("github") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": apperrors.ErrAuthMethodNotEnabled.Error()})
		}
		return next(c)
	}
}
