package secondary

import "aegis/internal/domain/entities"

type OAuthProviderRequests interface {
	ExchangeCodeForUserInfos(code, state, redirectUri string) (*entities.UserInfos, error)
}
