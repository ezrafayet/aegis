package ports

import (
	"net/http"
	"othnx/internal/domain/entities"
)

type UseCasesInterface interface {
	// For handlers
	GetSession(accessToken string) (domain.Session, error)
	Logout(refreshToken string) (*http.Cookie, *http.Cookie, error)

	// For middlewares
	CheckAndRefreshToken(accessToken, refreshToken string, forceRefresh bool) (*http.Cookie, *http.Cookie, error)
}

type OAuthUseCasesInterface interface {
	// For handlers
	GetAuthURL(redirectUri string) (string, error)
	ExchangeCode(code, state string) (*http.Cookie, *http.Cookie, error)

	// For middlewares
	CheckAuthEnabled(provider string) bool
}
