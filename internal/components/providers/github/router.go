package github

import (
	"othnx/internal/domain"

	"github.com/labstack/echo/v4"
)

type OAuthGithubRouter struct {
	Handlers       domain.OAuthProviderHandlers
	AuthMiddleware domain.OAuthMiddlewares
}

var _ domain.OAuthRouter = OAuthGithubRouter{}

func NewOAuthGithubRouter(h domain.OAuthProviderHandlers, m domain.OAuthMiddlewares) OAuthGithubRouter {
	return OAuthGithubRouter{
		Handlers:       h,
		AuthMiddleware: m,
	}
}

func (r OAuthGithubRouter) AttachRoutes(e *echo.Echo) {
	group := e.Group("/auth/github", r.AuthMiddleware.CheckAuthEnabled)
	group.GET("", r.Handlers.GetAuthURL)
	group.POST("/callback", r.Handlers.ExchangeCode)
}
