package domain

import (
	"othnx/pkg/apperrors"
	"time"
)

// todo: move this business logic somewhere
func GenerateTokensForUser(user User, deviceID string, config Config, refreshTokenRepository RefreshTokenRepository) (accessToken string, atExpiresAt int64, refreshToken string, rtExpiresAt int64, err error) {
	deviceFingerprint, err := GenerateDeviceFingerprint(deviceID)
	if err != nil {
		return "", -1, "", -1, err
	}

	err = refreshTokenRepository.DeleteRefreshTokenByDeviceFingerprint(user.ID, deviceFingerprint)
	if err != nil {
		return "", -1, "", -1, err
	}

	validRefreshTokens, err := refreshTokenRepository.CountValidRefreshTokensForUser(user.ID)
	if err != nil {
		return "", -1, "", -1, err
	}

	if validRefreshTokens >= 5 {
		return "", -1, "", -1, apperrors.ErrTooManyRefreshTokens
	}

	_ = refreshTokenRepository.CleanExpiredTokens(user.ID)

	newRefreshToken, rtExpiresAt, err := NewRefreshToken(user, deviceFingerprint, config)
	if err != nil {
		return "", -1, "", -1, err
	}
	err = refreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return "", -1, "", -1, err
	}

	accessToken, atExpiresAt, err = NewAccessToken(CustomClaims{
		UserID:       user.ID,
		EarlyAdopter: user.EarlyAdopter,
		Metadata:     user.Metadata,
		RolesValues:  user.RolesValues(),
	}, config, time.Now())
	if err != nil {
		return "", -1, "", -1, err
	}

	return accessToken, atExpiresAt, newRefreshToken.Token, rtExpiresAt, nil
}
