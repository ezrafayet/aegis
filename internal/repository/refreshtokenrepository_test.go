package repository

import (
	"othnx/internal/domain"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRefreshTokenRepository(t *testing.T) {
	t.Run("should create a refresh token", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		refreshToken, _, err := domain.NewRefreshToken(domain.User{ID: "123"}, "device-id", domain.Config{
			JWT: domain.JWTConfig{
				Secret: "xxxsecret",
				AccessTokenExpirationMin: 15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.CreateRefreshToken(refreshToken)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var retrievedToken domain.RefreshToken
		result := db.Model(&domain.RefreshToken{}).Where("token = ?", refreshToken.Token).First(&retrievedToken)
		if result.Error != nil {
			t.Fatal("expected no error", result.Error)
		}
		if retrievedToken.Token != refreshToken.Token {
			t.Fatal("expected token to be the same", retrievedToken.Token, refreshToken.Token)
		}
		if retrievedToken.UserID != refreshToken.UserID {
			t.Fatal("expected user_id to be the same", retrievedToken.UserID, refreshToken.UserID)
		}
		if retrievedToken.DeviceFingerprint != refreshToken.DeviceFingerprint {
			t.Fatal("expected device_fingerprint to be the same", retrievedToken.DeviceFingerprint, refreshToken.DeviceFingerprint)
		}
		if !retrievedToken.CreatedAt.Equal(refreshToken.CreatedAt) {
			t.Fatal("expected created_at to be the same", retrievedToken.CreatedAt, refreshToken.CreatedAt)
		}
		if !retrievedToken.ExpiresAt.Equal(refreshToken.ExpiresAt) {
			t.Fatal("expected expires_at to be the same", retrievedToken.ExpiresAt, refreshToken.ExpiresAt)
		}
	})
}
