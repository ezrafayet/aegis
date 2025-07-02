package repositories

import (
	"aegis/internal/domain/entities"
	"aegis/pkg/fingerprint"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateRefreshToken(t *testing.T) {
	t.Run("should create a refresh token", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		refreshToken, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, "device-id", entities.Config{
			JWT: entities.JWTConfig{
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
		var retrievedToken entities.RefreshToken
		result := db.Model(&entities.RefreshToken{}).Where("token = ?", refreshToken.Token).First(&retrievedToken)
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
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		deviceFingerprint1, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		deviceFingerprint2, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshToken1, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint1, entities.Config{
			JWT: entities.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshToken2, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint2, entities.Config{
			JWT: entities.JWTConfig{
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

func TestCountValidRefreshTokensForUser(t *testing.T) {
	t.Run("should count valid refresh tokens for user", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenExpired, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint, entities.Config{
			JWT: entities.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenExpired.ExpiresAt = time.Now().Add(-time.Hour * 24)
		err = refreshTokenRepository.CreateRefreshToken(refreshTokenExpired)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenActive, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint, entities.Config{
			JWT: entities.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.CreateRefreshToken(refreshTokenActive)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		count, err := refreshTokenRepository.CountValidRefreshTokensForUser("123")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		if count != 1 {
			t.Fatal("expected count to be 0", count)
		}
	})
}

func TestCleanExpiredTokens(t *testing.T) {
	t.Run("should clean expired tokens", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenExpired, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint, entities.Config{
			JWT: entities.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenExpired.ExpiresAt = time.Now().Add(-time.Hour * 24)
		err = refreshTokenRepository.CreateRefreshToken(refreshTokenExpired)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenActive, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint, entities.Config{
			JWT: entities.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.CreateRefreshToken(refreshTokenActive)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.CleanExpiredTokens("123")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var countActive int64
		result := db.Model(&entities.RefreshToken{}).Where("user_id = ? AND token = ?", "123", refreshTokenActive.Token).Count(&countActive)
		if result.Error != nil {
			t.Fatal("expected no error", result.Error)
		}
		if countActive != 1 {
			t.Fatal("expected count to be 1", countActive)
		}
		var countExpired int64
		result = db.Model(&entities.RefreshToken{}).Where("user_id = ? AND token = ?", "123", refreshTokenExpired.Token).Count(&countExpired)
		if result.Error != nil {
			t.Fatal("expected no error", result.Error)
		}
		if countExpired != 0 {
			t.Fatal("expected count to be 0", countExpired)
		}
	})
}

func TestDeleteRefreshToken(t *testing.T) {
	t.Run("should delete a refresh token", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
		refreshTokenRepository := NewRefreshTokenRepository(db)
		deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenExpired, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint, entities.Config{
			JWT: entities.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenExpired.ExpiresAt = time.Now().Add(-time.Hour * 24)
		err = refreshTokenRepository.CreateRefreshToken(refreshTokenExpired)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshTokenActive, _, err := entities.NewRefreshToken(entities.User{ID: "123"}, deviceFingerprint, entities.Config{
			JWT: entities.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.CreateRefreshToken(refreshTokenActive)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = refreshTokenRepository.DeleteRefreshToken(refreshTokenExpired.Token)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var count int64
		result := db.Model(&entities.RefreshToken{}).Where("user_id = ? AND token = ?", "123", refreshTokenActive.Token).Count(&count)
		if result.Error != nil {
			t.Fatal("expected no error", result.Error)
		}
		if count != 1 {
			t.Fatal("expected count to be 1", count)
		}
	})
}
