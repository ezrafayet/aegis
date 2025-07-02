package services

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"aegis/pkg/apperrors"
	"aegis/pkg/fingerprint"
	"aegis/pkg/jwtgen"
	"time"
)

type TokenService struct {
	refreshTokenRepository secondary.RefreshTokenRepository
	config                 entities.Config
}

func NewTokenService(refreshTokenRepository secondary.RefreshTokenRepository, config entities.Config) *TokenService {
	return &TokenService{
		refreshTokenRepository: refreshTokenRepository,
		config:                 config,
	}
}

// GenerateTokensForUser creates new access and refresh tokens for a user
// It handles device fingerprinting, token cleanup, and validation
func (s *TokenService) GenerateTokensForUser(user entities.User, deviceID string) (accessToken string, atExpiresAt int64, refreshToken string, rtExpiresAt int64, err error) {
	// Generate device fingerprint
	deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint(deviceID)
	if err != nil {
		return "", -1, "", -1, err
	}

	// Delete existing refresh token for this device
	err = s.refreshTokenRepository.DeleteRefreshTokenByDeviceFingerprint(user.ID, deviceFingerprint)
	if err != nil {
		return "", -1, "", -1, err
	}

	// Check token limit
	validRefreshTokens, err := s.refreshTokenRepository.CountValidRefreshTokensForUser(user.ID)
	if err != nil {
		return "", -1, "", -1, err
	}

	if validRefreshTokens >= 5 {
		return "", -1, "", -1, apperrors.ErrTooManyRefreshTokens
	}

	// Clean expired tokens
	_ = s.refreshTokenRepository.CleanExpiredTokens(user.ID)

	// Create new refresh token
	newRefreshToken, rtExpiresAt, err := entities.NewRefreshToken(user, deviceFingerprint, s.config)
	if err != nil {
		return "", -1, "", -1, err
	}
	err = s.refreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return "", -1, "", -1, err
	}

	// Generate access token
	accessToken, atExpiresAt, err = jwtgen.Generate(entities.CustomClaims{
		UserID:       user.ID,
		EarlyAdopter: user.EarlyAdopter,
		Metadata:     user.Metadata,
		Roles:        user.RolesValues(),
	}, s.config, time.Now())
	if err != nil {
		return "", -1, "", -1, err
	}

	return accessToken, atExpiresAt, newRefreshToken.Token, rtExpiresAt, nil
}
