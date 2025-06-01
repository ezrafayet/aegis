package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID            string    `json:"id" gorm:"type:uuid"`
	CreatedAt         time.Time `json:"created_at" gorm:"index;not null"`
	ExpiresAt         time.Time `json:"expires_at" gorm:"index;not null"`
	Token             string    `json:"token" gorm:"primaryKey;type:varchar(150);not null"`
	// DeviceFingerprint string   `json:"device_fingerprint" gorm:"type:varchar(150)"`
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
	GetValidRefreshTokensByUserID(userID string) ([]RefreshToken, error)
	CleanExpiredTokens(userID string) error
	// GetRefreshTokensByUserID(userID string) ([]RefreshToken, error)
	// DeleteRefreshToken(token string) error
}
