package secondary

import "aegis/internal/domain/entities"

type OAuthProviderConfig interface {
	IsEnabled() (bool, error)
	GetName() string
	GetOauthRedirectURL(redirectUrl, state string) (string, error)
}

type OAuthProviderRequests interface {
	ExchangeCodeForUserInfos(code, state, redirectUri string) (*entities.UserInfos, error)
}

type OAuthProviderInterface interface {
	OAuthProviderConfig
	OAuthProviderRequests
}
