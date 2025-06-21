package repositories

import "othnx/internal/core/domain"

type RefreshTokenRepository interface {
	CreateRefreshToken(refreshToken domain.RefreshToken) error
	GetRefreshTokenByToken(token string) (domain.RefreshToken, error)
	CountValidRefreshTokensForUser(userID string) (int, error)
	CleanExpiredTokens(userID string) error
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByDeviceFingerprint(userID, deviceFingerprint string) error
}
