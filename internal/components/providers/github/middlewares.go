package github

import (
	"aegix/internal/components/providers"
	"aegix/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OAuthGithubMiddlewares struct {
	Config domain.Config
}

var _ providers.OAuthMiddlewares = OAuthGithubMiddlewares{}

func NewOAuthGithubMiddlewares(c domain.Config) OAuthGithubMiddlewares {
	return OAuthGithubMiddlewares{
		Config: c,
	}
}

func (m OAuthGithubMiddlewares) CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !m.Config.Auth.Providers.GitHub.Enabled {
			return c.JSON(http.StatusForbidden, map[string]string{"error": domain.ErrAuthMethodNotEnabled.Error()})
		}
		return next(c)
	}
}
