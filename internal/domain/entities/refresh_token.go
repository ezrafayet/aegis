package entities

import (
	"othnx/internal/infrastructure/config"
	"othnx/pkg/tokengen"
	"time"
)

type RefreshToken struct {
	UserID            string    `json:"id" gorm:"type:uuid"`
	CreatedAt         time.Time `json:"created_at" gorm:"index;not null"`
	ExpiresAt         time.Time `json:"expires_at" gorm:"index;not null"`
	Token             string    `json:"token" gorm:"primaryKey;type:char(32);not null"`
	DeviceFingerprint string    `json:"device_fingerprint" gorm:"type:char(32);index;not null"`
	// relations
	User   User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

// todo: make tokens unique per deviceID
// todo: cascade delete tokens on user deletion

func (r RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}

func NewRefreshToken(user User, deviceFingerprint string, config config.Config) (RefreshToken, int64, error) {
	token, err := tokengen.Generate("refresh_", 12)
	if err != nil {
		return RefreshToken{}, -1, err
	}
	createdAt := time.Now()
	expiresAt := createdAt.AddDate(0, 0, config.JWT.RefreshTokenExpirationDays)
	return RefreshToken{
		UserID:            user.ID,
		CreatedAt:         createdAt,
		ExpiresAt:         expiresAt,
		Token:             token,
		DeviceFingerprint: deviceFingerprint,
	}, expiresAt.Unix(), nil
}
