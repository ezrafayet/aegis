package github

import (
	"aegix/internal/domain"
	"aegix/internal/providers"
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
			return c.JSON(http.StatusForbidden, map[string]string{"error": providers.AuthMethodNotEnabled.Error()})
		}
		return next(c)
	}
}
