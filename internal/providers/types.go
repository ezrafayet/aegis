package providers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrNoUser               = errors.New("no_user")
	ErrUserBlocked          = errors.New("user_blocked")
	ErrUserDeleted          = errors.New("user_deleted")
	ErrWrongAuthMethod      = errors.New("wrong_auth_method")
	ErrNoRefreshToken       = errors.New("no_refresh_token")
	ErrTooManyRefreshTokens = errors.New("too_many_refresh_tokens")
)

var (
	ErrAuthMethodNotEnabled = errors.New("auth_method_not_enabled")
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
