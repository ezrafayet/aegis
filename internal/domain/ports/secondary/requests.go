package secondary

import "othnx/internal/domain/entities"

type OAuthProviderRequests interface {
	GetUserInfos(code, state, redirectUri string) (*entities.UserInfos, error)
}
