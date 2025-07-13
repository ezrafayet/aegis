package providers

type OAuthProviderConfig interface {
	IsEnabled() bool
	GetName() string
	GetOauthRedirectURL(state string) string
}

type OAuthProviderRequests interface {
	ExchangeCodeForUserInfos(code, state string) (*UserInfos, error)
}

type OAuthProviderInterface interface {
	OAuthProviderConfig
	OAuthProviderRequests
}

type OAuthRepository struct {
	Name         string
	Enabled      bool
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type UserInfos struct {
	Name   string
	Email  string
	Avatar string
}
