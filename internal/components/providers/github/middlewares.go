package github

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"othnx/internal/components/providers/providersports"
	"othnx/internal/domain"
	"othnx/pkg/apperrors"
)

type OAuthGithubMiddlewares struct {
	Config domain.Config
}

var _ providersports.OAuthMiddlewares = OAuthGithubMiddlewares{}

func NewOAuthGithubMiddlewares(c domain.Config) OAuthGithubMiddlewares {
	return OAuthGithubMiddlewares{
		Config: c,
	}
}

func (m OAuthGithubMiddlewares) CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !m.Config.Auth.Providers.GitHub.Enabled {
			return c.JSON(http.StatusForbidden, map[string]string{"error": apperrors.ErrAuthMethodNotEnabled.Error()})
		}
		return next(c)
	}
}
