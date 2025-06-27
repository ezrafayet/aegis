package middlewares

import (
	"net/http"
	"othnx/internal/domain/ports/primary_ports"
	"othnx/internal/infrastructure/config"
	"othnx/pkg/apperrors"

	"github.com/labstack/echo/v4"
)

// some factory

type OAuthMiddlewaresInterface interface {
	CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc
}

type OAuthGithubMiddlewares struct {
	Config  config.Config
	Service primaryports.OAuthUseCasesInterface
}

var _ OAuthMiddlewaresInterface = OAuthGithubMiddlewares{}

func NewOAuthGithubMiddlewares(c config.Config) OAuthGithubMiddlewares {
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
