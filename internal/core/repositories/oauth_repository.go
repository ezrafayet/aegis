package repositories

import "othnx/internal/core/domain"

type OAuthProviderRepository interface {
	GetUserInfos(code, state, redirectUri string) (*domain.UserInfos, error)
}
