package domain

import (
	"time"
)

type RefreshToken struct {
	UserID    string    `json:"id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"index;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null"`
	Token     string    `json:"token" gorm:"primaryKey;type:char(32);not null"`
	// DeviceFingerprint string   `json:"device_fingerprint" gorm:"type:varchar(150)"`
}

func (r RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}

func NewRefreshToken(user User, config Config) (RefreshToken, int64, error) {
	token, err := GenerateRandomToken("refresh_", 12)
	if err != nil {
		return RefreshToken{}, -1, err
	}
	createdAt := time.Now()
	expiresAt := createdAt.AddDate(0, 0, config.JWT.RefreshTokenExpirationDays)
	return RefreshToken{
		UserID:    user.ID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
		Token:     token,
	}, expiresAt.Unix(), nil
}

type RefreshTokenRepository interface {
	CreateRefreshToken(refreshToken RefreshToken) error
	GetRefreshTokenByToken(token string) (RefreshToken, error)
	GetValidRefreshTokensByUserID(userID string) ([]RefreshToken, error)
	CleanExpiredTokens(userID string) error
	DeleteRefreshToken(token string) error
}
