package secondary

import "othnx/internal/domain/entities"

type OAuthProviderRequests interface {
	ExchangeCodeForUserInfos(code, state, redirectUri string) (*entities.UserInfos, error)
}
