package github

import (
	"aegix/internal/components/providers/providersports"

	"github.com/labstack/echo/v4"
)

type OAuthGithubRouter struct {
	Handlers       providersports.OAuthProviderHandlers
	AuthMiddleware providersports.OAuthMiddlewares
}

var _ providersports.OAuthRouter = OAuthGithubRouter{}

func NewOAuthGithubRouter(h providersports.OAuthProviderHandlers, m providersports.OAuthMiddlewares) OAuthGithubRouter {
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
