package domain

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type RefreshToken struct {
	UserID    string    `json:"id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"index;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null"`
	Token     string    `json:"token" gorm:"primaryKey;type:char(32);not null"`
	DeviceFingerprint  string    `json:"device_fingerprint" gorm:"type:char(32);index;not null"`
}

func (r RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}

func NewRefreshToken(user User, deviceID string, config Config) (RefreshToken, int64, error) {
	token, err := GenerateRandomToken("refresh_", 12)
	if err != nil {
		return RefreshToken{}, -1, err
	}
	createdAt := time.Now()
	expiresAt := createdAt.AddDate(0, 0, config.JWT.RefreshTokenExpirationDays)
	deviceFingerprint, err := GenerateDeviceFingerprint(deviceID)
	if err != nil {
		return RefreshToken{}, -1, err
	}
	return RefreshToken{
		UserID:    user.ID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
		Token:     token,
		DeviceFingerprint: deviceFingerprint,
	}, expiresAt.Unix(), nil
}

type RefreshTokenRepository interface {
	CreateRefreshToken(refreshToken RefreshToken) error
	GetRefreshTokenByToken(token string) (RefreshToken, error)
	GetValidRefreshTokensByUserID(userID string) ([]RefreshToken, error)
	CleanExpiredTokens(userID string) error
	DeleteRefreshToken(token string) error
}

func GenerateDeviceFingerprint(deviceID string) (string, error) {
	trimmed := strings.TrimSpace(deviceID)
	if trimmed == "" {
		trimmed = "default-device-id"
	}
	normalized := strings.ToLower(trimmed)
	transformer := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, err := transform.String(transformer, normalized)
	if err != nil {
		result = normalized
	}
	result = strings.Join(strings.Fields(result), " ")
	hash := md5.Sum([]byte(result))
	return hex.EncodeToString(hash[:]), nil
}
