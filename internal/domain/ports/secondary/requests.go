package secondary

import "aegis/internal/domain/entities"

type OAuthProviderRequests interface {
	GetName() string
	ExchangeCodeForUserInfos(code, state, redirectUri string) (*entities.UserInfos, error)
}
