package github

import (
	"othnx/internal/domain"
	"othnx/internal/repository"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockRepository struct{}

func (p MockRepository) GetUserInfos(code, state, redirectUri string) (*domain.UserInfos, error) {
	userInfos := domain.UserInfos{
		Name:   "some-name",
		Email:  "some-email",
		Avatar: "some-avatar",
	}
	return &userInfos, nil
}

func TestOAuthGithubService_ExchangeCode(t *testing.T) {
	baseConfig := domain.Config{
		JWT: domain.JWTConfig{
			Secret:                     "some-secret",
			AccessTokenExpirationMin:   1,
			RefreshTokenExpirationDays: 1,
		},
	}
	prepare := func(t *testing.T) (OAuthGithubService, repository.UserRepository, repository.RefreshTokenRepository, *gorm.DB) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}, &domain.State{})
		refreshTokenRepository := repository.NewRefreshTokenRepository(db)
		oauthProviderRepository := MockRepository{}
		userRepository := repository.NewUserRepository(db)
		stateRepository := repository.NewStateRepository(db)
		authService := NewOAuthGithubService(baseConfig, oauthProviderRepository, &userRepository, &refreshTokenRepository, &stateRepository)
		return authService, userRepository, refreshTokenRepository, db
	}
	t.Run("should create a user if no user exists", func(t *testing.T) {
		ghService, _, _, db := prepare(t)
		db.Create(&domain.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
		_, _, err := ghService.ExchangeCode("some-code", "some-state")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var user domain.User
		db.First(&user, "email = ?", "some-email")
		if user.Email != "some-email" {
			t.Fatal("expected user to be created", user.Email)
		}
	})
	t.Run("should not create a user if user already exists", func(t *testing.T) {
		ghService, _, _, db := prepare(t)
		db.Create(&domain.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
		newUser, err := domain.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		_, _, err = ghService.ExchangeCode("some-code", "some-state")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var count int64
		db.Model(&domain.User{}).Where("email = ?", "some-email").Count(&count)
		if count != 1 {
			t.Fatal("expected user to be created", count)
		}
	})
	t.Run("should not return tokens if user is blocked", func(t *testing.T) {
		ghService, _, _, db := prepare(t)
		db.Create(&domain.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
		newUser, err := domain.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		now := time.Now()
		newUser.BlockedAt = &now
		db.Create(&newUser)
		_, _, err = ghService.ExchangeCode("some-code", "some-state")
		if err.Error() != domain.ErrUserBlocked.Error() {
			t.Fatal("expected error ErrUserBlocked", err)
		}
	})
	t.Run("should not return tokens if user is deleted", func(t *testing.T) {
		ghService, _, _, db := prepare(t)
		db.Create(&domain.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
		newUser, err := domain.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		now := time.Now()
		newUser.DeletedAt = &now
		db.Create(&newUser)
		_, _, err = ghService.ExchangeCode("some-code", "some-state")
		if err.Error() != domain.ErrUserDeleted.Error() {
			t.Fatal("expected error ErrUserDeleted", err)
		}
	})
	t.Run("should issue tokens if user does not have a valid refresh token", func(t *testing.T) {
		ghService, _, _, db := prepare(t)
		db.Create(&domain.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
		at, rt, err := ghService.ExchangeCode("some-code", "some-state")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var user domain.User
		db.First(&user, "email = ?", "some-email")
		if user.Email != "some-email" {
			t.Fatal("expected user to be created", user.Email)
		}
		if at == nil || at.Value == "" {
			t.Fatal("expected access token to be nil", at)
		}
		if rt == nil || rt.Value == "" {
			t.Fatal("expected refresh token to be nil", rt)
		}
		var refreshToken domain.RefreshToken
		db.First(&refreshToken, "user_id = ?", user.ID)
		if refreshToken.Token == "" {
			t.Fatal("expected refresh token to be created", refreshToken.Token)
		}
	})
	t.Run("should issue new tokens if user already has a valid refresh token, and delete the previous one", func(t *testing.T) {
		ghService, _, _, db := prepare(t)
		db.Create(&domain.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
		_, rt1, err := ghService.ExchangeCode("some-code", "some-state")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		var user domain.User
		db.First(&user, "email = ?", "some-email")
		if user.Email != "some-email" {
			t.Fatal("expected user to be created", user.Email)
		}
		if rt1 == nil || rt1.Value == "" {
			t.Fatal("expected refresh token to be nil", rt1)
		}
		db.Create(&domain.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
		_, rt2, err := ghService.ExchangeCode("some-code", "some-state")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		if rt2 == nil || rt2.Value == "" {
			t.Fatal("expected refresh token to be nil", rt2)
		}
		var count int64
		db.Model(&domain.RefreshToken{}).Where("user_id = ?", user.ID).Count(&count)
		if count != 1 {
			t.Fatal("expected refresh token to be deleted", count)
		}
		var refreshToken domain.RefreshToken
		db.First(&refreshToken, "user_id = ?", user.ID)
		if refreshToken.Token == rt1.Value {
			t.Fatal("expected refresh token to be deleted", refreshToken.Token)
		}
		if refreshToken.Token != rt2.Value {
			t.Fatal("expected refresh token to be created", refreshToken.Token)
		}
	})
	// todo: add specific tests when device id is handled
}
