package repository

import (
	"fmt"
	"othnx/internal/domain"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateRefreshToken(t *testing.T) {
	t.Run("should create a refresh token", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		refreshToken, _, err := domain.NewRefreshToken(domain.User{ID: "123"}, "device-id", domain.Config{
			JWT: domain.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
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

func TestGetRefreshTokenByToken(t *testing.T) {
	t.Run("should get a refresh token by token", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		deviceFingerprint1, err := domain.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		deviceFingerprint2, err := domain.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshToken1, _, err := domain.NewRefreshToken(domain.User{ID: "123"}, deviceFingerprint1, domain.Config{
			JWT: domain.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshToken2, _, err := domain.NewRefreshToken(domain.User{ID: "123"}, deviceFingerprint2, domain.Config{
			JWT: domain.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.CreateRefreshToken(refreshToken1)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.CreateRefreshToken(refreshToken2)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		retrievedToken, err := refreshTokenRepository.GetRefreshTokenByToken(refreshToken2.Token)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		fmt.Println("retrievedToken", retrievedToken.Token)
		fmt.Println("refreshToken2", refreshToken2.Token)
		fmt.Println("refreshToken1", refreshToken1.Token)
		fmt.Println("refreshToken2", refreshToken2.Token)
		if retrievedToken.Token != refreshToken2.Token {
			t.Fatal("expected token to be the same", retrievedToken.Token, refreshToken2.Token)
		}
	})
}
