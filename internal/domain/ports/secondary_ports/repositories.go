package secondary_ports

import "othnx/internal/domain/entities"

type OAuthProviderRepository interface {
	GetUserInfos(code, state, redirectUri string) (*entities.UserInfos, error)
}

type RefreshTokenRepository interface {
	CreateRefreshToken(refreshToken entities.RefreshToken) error
	GetRefreshTokenByToken(token string) (entities.RefreshToken, error)
	CountValidRefreshTokensForUser(userID string) (int, error)
	CleanExpiredTokens(userID string) error
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByDeviceFingerprint(userID, deviceFingerprint string) error
}

type StateRepository interface {
	CreateState(state entities.State) error
	GetAndDeleteState(value string) (entities.State, error)
}

type UserRepository interface {
	CreateUser(user entities.User, roles []entities.Role) error
	GetUserByID(userID string) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	DoesNameExist(nameFingerprint string) (bool, error)
	// add role
	// remove role
}
