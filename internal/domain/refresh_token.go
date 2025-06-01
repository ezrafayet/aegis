package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID    string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
}

func (r RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}

func NewRefreshToken(userID string, validityInDays int) RefreshToken {
	createdAt := time.Now()
	expiresAt := createdAt.AddDate(0, 0, validityInDays)
	return RefreshToken{
		UserID:    userID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
		Token:     uuid.New().String(),
	}
}

type RefreshTokenRepository interface {
	CreateRefreshToken(refreshToken RefreshToken) error
	GetRefreshTokenByToken(token string) (RefreshToken, error)
	// GetRefreshTokensByUserID(userID string) ([]RefreshToken, error)
	// DeleteRefreshToken(token string) error
}
