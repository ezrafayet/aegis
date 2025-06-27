package usecases

import (
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/secondary_ports"
	"othnx/internal/infrastructure/config"
	"othnx/internal/infrastructure/repositories"
	"othnx/pkg/apperrors"
	"othnx/pkg/jwtgen"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCheckAndRefreshToken(t *testing.T) {
	baseConfig := config.Config{
		JWT: config.JWTConfig{
			Secret:                     "some-secret",
			AccessTokenExpirationMin:   1,
			RefreshTokenExpirationDays: 1,
		},
	}
	prepare := func(t *testing.T) (AuthService, secondaryports.UserRepository, secondaryports.RefreshTokenRepository, *gorm.DB) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{})
		refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
		userRepository := repositories.NewUserRepository(db)
		authService := NewService(baseConfig, &refreshTokenRepository, &userRepository)
		return authService, &userRepository, &refreshTokenRepository, db
	}
	t.Run("invalid access token gets rejected", func(t *testing.T) {
		authService, _, _, _ := prepare(t)
		_, err := authService.CheckAndRefreshToken("invalid", "invalid", false)
		if err.Error() != apperrors.ErrAccessTokenInvalid.Error() {
			t.Fatal("expected error ErrInvalidAccessToken", err)
		}
	})
	t.Run("expired access token and expired refresh token gets rejected", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		refreshToken, _, err := entities.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshToken.ExpiresAt = time.Now().Add(-time.Hour * 24)
		db.Save(&refreshToken)
		accessToken, _, err := jwtgen.Generate(entities.CustomClaims{UserID: newUser.ID}, config.Config{}, time.Now().Add(-time.Hour*24))
		if err != nil {
			t.Fatal("expected no error", err)
		}
		_, err = authService.CheckAndRefreshToken(accessToken, refreshToken.Token, false)
		if err.Error() != apperrors.ErrRefreshTokenExpired.Error() {
			t.Fatal("expected error ErrRefreshTokenExpired", err)
		}
	})
	t.Run("expired access token and valid refresh token passes, new tokens are generated", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		refreshToken, _, err := entities.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Save(&refreshToken)
		accessToken, _, err := jwtgen.Generate(entities.CustomClaims{UserID: newUser.ID}, config.Config{}, time.Now().Add(-time.Hour*24))
		if err != nil {
			t.Fatal("expected no error", err)
		}
		at, err := authService.CheckAndRefreshToken(accessToken, refreshToken.Token, false)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var retrievedRefreshToken entities.RefreshToken
		db.First(&retrievedRefreshToken, "user_id = ?", newUser.ID)
		if retrievedRefreshToken.Token == refreshToken.Token {
			t.Fatal("expected refresh token to be different")
		}
		if at.AccessToken == accessToken {
			t.Fatal("expected access token to be different")
		}
	})
	t.Run("valid access token does not trigger refresh", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		refreshToken, _, err := entities.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Save(&refreshToken)
		accessToken, _, err := jwtgen.Generate(entities.CustomClaims{UserID: newUser.ID}, config.Config{}, time.Now())
		if err != nil {
			t.Fatal("expected no error", err)
		}
		_, err = authService.CheckAndRefreshToken(accessToken, refreshToken.Token, false)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var retrievedRefreshToken entities.RefreshToken
		db.First(&retrievedRefreshToken, "user_id = ?", newUser.ID)
		if retrievedRefreshToken.Token != refreshToken.Token {
			t.Fatal("expected refresh token to be the same")
		}
	})

	t.Run("valid access token does trigger refresh if forceRefresh is true", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		refreshToken, _, err := entities.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Save(&refreshToken)
		accessToken, _, err := jwtgen.Generate(entities.CustomClaims{UserID: newUser.ID}, config.Config{}, time.Now())
		if err != nil {
			t.Fatal("expected no error", err)
		}
		_, err = authService.CheckAndRefreshToken(accessToken, refreshToken.Token, true)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var retrievedRefreshToken entities.RefreshToken
		db.First(&retrievedRefreshToken, "user_id = ?", newUser.ID)
		if retrievedRefreshToken.Token == refreshToken.Token {
			t.Fatal("expected refresh token to be different")
		}
	})
}
