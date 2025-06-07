package github

import (
	"aegix/internal/components/providers"

	"github.com/labstack/echo/v4"
)

type OAuthGithubRouter struct {
	Handlers       providers.OAuthProviderHandlers
	AuthMiddleware providers.OAuthMiddlewares
}

var _ providers.OAuthRouter = OAuthGithubRouter{}

func NewOAuthGithubRouter(h providers.OAuthProviderHandlers, m providers.OAuthMiddlewares) OAuthGithubRouter {
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
