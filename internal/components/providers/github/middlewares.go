package github

import (
	"net/http"
	"othnx/internal/domain"
	"othnx/pkg/apperrors"

	"github.com/labstack/echo/v4"
)

type OAuthGithubMiddlewares struct {
	Config domain.Config
}

var _ domain.OAuthMiddlewares = OAuthGithubMiddlewares{}

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
