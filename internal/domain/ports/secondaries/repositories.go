package ports

import "othnx/internal/domain/entities"

type OAuthProviderRepository interface {
	GetUserInfos(code, state, redirectUri string) (*domain.UserInfos, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(refreshToken domain.RefreshToken) error
	GetRefreshTokenByToken(token string) (domain.RefreshToken, error)
	CountValidRefreshTokensForUser(userID string) (int, error)
	CleanExpiredTokens(userID string) error
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByDeviceFingerprint(userID, deviceFingerprint string) error
}

type StateRepository interface {
	CreateState(state domain.State) error
	GetAndDeleteState(value string) (domain.State, error)
}

type UserRepository interface {
	CreateUser(user domain.User, roles []domain.Role) error
	GetUserByID(userID string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	DoesNameExist(nameFingerprint string) (bool, error)
	// add role
	// remove role
}
