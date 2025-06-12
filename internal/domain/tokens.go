package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateRandomToken(prefix string, nPairs int) (string, error) {
	bytes := make([]byte, nPairs)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(bytes), nil
}

// todo: move this business logic somewhere
func GenerateTokensForUser(user User, deviceID string, config Config, refreshTokenRepository RefreshTokenRepository) (accessToken string, atExpiresAt int64, refreshToken string, rtExpiresAt int64, err error) {
	validRefreshTokens, err := refreshTokenRepository.GetValidRefreshTokensByUserID(user.ID)
	if err != nil {
		return "", -1, "", -1, err
	}
	if len(validRefreshTokens) > 10 {
		return "", -1, "", -1, ErrTooManyRefreshTokens
	}

	_ = refreshTokenRepository.CleanExpiredTokens(user.ID)

	// device-id will need to be injected
	newRefreshToken, rtExpiresAt, err := NewRefreshToken(user, deviceID, config)
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
	}, config, time.Now())
	if err != nil {
		return "", -1, "", -1, err
	}

	return accessToken, atExpiresAt, newRefreshToken.Token, rtExpiresAt, nil
}
