package auth

import (
	"othnx/internal/domain"
	"othnx/internal/repository"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCheckAndRefreshToken(t *testing.T) {
	prepare := func(t *testing.T) (AuthService, repository.UserRepository, repository.RefreshTokenRepository, *gorm.DB) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{}, &domain.RefreshToken{})
		refreshTokenRepository := repository.NewRefreshTokenRepository(db)
		userRepository := repository.NewUserRepository(db)
		authService := NewAuthService(domain.Config{}, refreshTokenRepository, userRepository)
		return authService, userRepository, refreshTokenRepository, db
	}
	t.Run("invalid access token gets rejected", func(t *testing.T) {
		authService, _, _, _ := prepare(t)
		_, _, _, err := authService.CheckAndRefreshToken("invalid", "invalid")
		if err.Error() != domain.ErrInvalidAccessToken.Error() {
			t.Fatal("expected error", err)
		}
	})
	t.Run("expired access token and expired refresh token gets rejected", func(t *testing.T) {})
	t.Run("expired access token and valid refresh token passes, new tokens are generated", func(t *testing.T) {})
	t.Run("valid access token does not trigger refresh", func(t *testing.T) {})
}
