package auth

import (
	"othnx/internal/domain"
	"othnx/internal/repository"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCheckAndRefreshToken(t *testing.T) {
	baseConfig := domain.Config{
		JWT: domain.JWTConfig{
			Secret:                     "some-secret",
			AccessTokenExpirationMin:   1,
			RefreshTokenExpirationDays: 1,
		},
	}
	prepare := func(t *testing.T) (AuthService, repository.UserRepository, repository.RefreshTokenRepository, *gorm.DB) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})
		refreshTokenRepository := repository.NewRefreshTokenRepository(db)
		userRepository := repository.NewUserRepository(db)
		authService := NewAuthService(baseConfig, refreshTokenRepository, userRepository)
		return authService, userRepository, refreshTokenRepository, db
	}
	t.Run("invalid access token gets rejected", func(t *testing.T) {
		authService, _, _, _ := prepare(t)
		_, _, err := authService.CheckAndRefreshToken("invalid", "invalid")
		if err.Error() != domain.ErrInvalidAccessToken.Error() {
			t.Fatal("expected error ErrInvalidAccessToken", err)
		}
	})
	t.Run("expired access token and expired refresh token gets rejected", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := domain.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		refreshToken, _, err := domain.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshToken.ExpiresAt = time.Now().Add(-time.Hour * 24)
		db.Save(&refreshToken)
		accessToken, _, err := domain.NewAccessToken(domain.CustomClaims{UserID: newUser.ID}, domain.Config{}, time.Now().Add(-time.Hour*24))
		if err != nil {
			t.Fatal("expected no error", err)
		}
		_, _, err = authService.CheckAndRefreshToken(accessToken, refreshToken.Token)
		if err.Error() != domain.ErrRefreshTokenExpired.Error() {
			t.Fatal("expected error ErrRefreshTokenExpired", err)
		}
	})
	t.Run("expired access token and valid refresh token passes, new tokens are generated", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := domain.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		refreshToken, _, err := domain.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Save(&refreshToken)
		accessToken, _, err := domain.NewAccessToken(domain.CustomClaims{UserID: newUser.ID}, baseConfig, time.Now().Add(-time.Hour*24))
		if err != nil {
			t.Fatal("expected no error", err)
		}
		at, _, err := authService.CheckAndRefreshToken(accessToken, refreshToken.Token)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var retrievedRefreshToken domain.RefreshToken
		db.First(&retrievedRefreshToken, "user_id = ?", newUser.ID)
		if retrievedRefreshToken.Token == refreshToken.Token {
			t.Fatal("expected refresh token to be different")
		}
		if at.Value == accessToken {
			t.Fatal("expected access token to be different")
		}
	})
	t.Run("valid access token does not trigger refresh", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := domain.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		refreshToken, _, err := domain.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Save(&refreshToken)
		accessToken, _, err := domain.NewAccessToken(domain.CustomClaims{UserID: newUser.ID}, baseConfig, time.Now())
		if err != nil {
			t.Fatal("expected no error", err)
		}
		_, _, err = authService.CheckAndRefreshToken(accessToken, refreshToken.Token)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var retrievedRefreshToken domain.RefreshToken
		db.First(&retrievedRefreshToken, "user_id = ?", newUser.ID)
		if retrievedRefreshToken.Token != refreshToken.Token {
			t.Fatal("expected refresh token to be the same")
		}
	})
}
