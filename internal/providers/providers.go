package providers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type OAuthRouter interface {
	AttachRoutes(e *echo.Echo)
}

type OAuthProviderHandlers interface {
	GetAuthURL(c echo.Context) error
	ExchangeCode(c echo.Context) error
}

type OAuthMiddlewares interface {
	CheckAuthEnabled(next echo.HandlerFunc) echo.HandlerFunc
}

type OAuthProviderService interface {
	GetAuthURL(redirectUri string) (string, error)
	ExchangeCode(code, state string) (http.Cookie, http.Cookie, error) // returns a cookie
}

type OAuthProvider interface {
	GetUserInfos(code, state, redirectUri string) (*OAuthUser, error)
}

type OAuthUser struct {
	Name   string
	Email  string
	Avatar string
}
