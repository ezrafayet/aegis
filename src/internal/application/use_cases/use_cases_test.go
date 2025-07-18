package usecases

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"aegis/internal/infrastructure/repositories"
	"aegis/pkg/apperrors"
	"aegis/pkg/jwtgen"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCheckAndRefreshToken(t *testing.T) {
	baseConfig := entities.Config{
		JWT: entities.JWTConfig{
			Secret:                     "some-secret",
			AccessTokenExpirationMin:   1,
			RefreshTokenExpirationDays: 1,
		},
	}
	prepare := func(t *testing.T) (*UseCases, secondary.UserRepository, secondary.RefreshTokenRepository, *gorm.DB) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{}, &entities.Role{})
		refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
		userRepository := repositories.NewUserRepository(db)
		authService := NewService(baseConfig, refreshTokenRepository, userRepository)
		return authService, userRepository, refreshTokenRepository, db
	}
	t.Run("invalid access token gets rejected", func(t *testing.T) {
		authService, _, _, _ := prepare(t)
		_, err := authService.CheckAndRefreshToken("invalid", "invalid", false)
		if err.Error() != apperrors.ErrRefreshTokenInvalid.Error() {
			t.Fatal("expected error ErrRefreshTokenInvalid", err)
		}
	})
	t.Run("expired access token and expired refresh token gets rejected", func(t *testing.T) {
		authService, _, _, db := prepare(t)
		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		db.Create(&newUser)
		db.Save(&newUser)
		refreshToken, _, err := entities.NewRefreshToken(newUser, "some-device-id", baseConfig)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		refreshToken.ExpiresAt = time.Now().Add(-time.Hour * 24)
		db.Save(&refreshToken)
		newUser.Roles = []entities.Role{
			{UserID: newUser.ID, Value: "user"},
		}
		db.Save(&newUser)
		cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		accessToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now().Add(-time.Hour*24), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
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
		newUser.Roles = []entities.Role{
			{UserID: newUser.ID, Value: "user"},
		}
		db.Save(&newUser)
		cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		accessToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now().Add(-time.Hour*24), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
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
		newUser.Roles = []entities.Role{
			{UserID: newUser.ID, Value: "user"},
		}
		db.Save(&newUser)
		cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		accessToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now(), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
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
		newUser.Roles = []entities.Role{
			{UserID: newUser.ID, Value: "user"},
		}
		db.Save(&newUser)
		cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		accessToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now(), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
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

func TestAuthorize(t *testing.T) {
	baseConfig := entities.Config{
		JWT: entities.JWTConfig{
			Secret:                     "some-secret",
			AccessTokenExpirationMin:   1,
			RefreshTokenExpirationDays: 1,
		},
	}
	prepare := func(t *testing.T) (*UseCases, secondary.UserRepository, secondary.RefreshTokenRepository, *gorm.DB) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{}, &entities.Role{})
		refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
		userRepository := repositories.NewUserRepository(db)
		authService := NewService(baseConfig, refreshTokenRepository, userRepository)
		return authService, userRepository, refreshTokenRepository, db
	}

	t.Run("token", func(t *testing.T) {
		authService, _, _, _ := prepare(t)

		t.Run("no token returns an error", func(t *testing.T) {
			err := authService.Authorize("", []string{"user"})
			if err == nil {
				t.Fatal("expected error for empty token")
			}
			if err.Error() != apperrors.ErrAccessTokenInvalid.Error() {
				t.Fatal("expected error for invalid token", err.Error())
			}
		})

		t.Run("expired token returns an error", func(t *testing.T) {
			// Create a user and generate an expired token
			newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
			if err != nil {
				t.Fatal("expected no error", err)
			}
			newUser.Roles = []entities.Role{
				{UserID: newUser.ID, Value: "user"},
			}
			cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
			if err != nil {
				t.Fatal("expected no error", err)
			}
			// Generate token with past expiration
			expiredToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now().Add(-time.Hour*24), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
			if err != nil {
				t.Fatal("expected no error", err)
			}
			err = authService.Authorize(expiredToken, []string{"user"})
			if err == nil {
				t.Fatal("expected error for expired token")
			}
			if err.Error() != apperrors.ErrAccessTokenExpired.Error() {
				t.Fatal("expected error for expired token", err.Error())
			}
		})

		t.Run("malformed token returns an error", func(t *testing.T) {
			err := authService.Authorize("invalid.token.here", []string{"user"})
			if err == nil {
				t.Fatal("expected error for malformed token")
			}
			if err.Error() != apperrors.ErrAccessTokenInvalid.Error() {
				t.Fatal("expected error for malformed token", err.Error())
			}
		})
	})

	t.Run("roles", func(t *testing.T) {
		authService, _, _, _ := prepare(t)

		t.Run("no roles returns error", func(t *testing.T) {
			newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
			if err != nil {
				t.Fatal("expected no error", err)
			}
			newUser.Roles = []entities.Role{
				{UserID: newUser.ID, Value: "user"},
			}
			cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
			if err != nil {
				t.Fatal("expected no error", err)
			}
			validToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now(), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
			if err != nil {
				t.Fatal("expected no error", err)
			}
			err = authService.Authorize(validToken, []string{})
			if err.Error() != apperrors.ErrNoRoles.Error() {
				t.Fatal("expected error ErrNoRoles", err)
			}
		})

		t.Run("any role returns true even if no roles are present on user", func(t *testing.T) {
			newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
			if err != nil {
				t.Fatal("expected no error", err)
			}
			newUser.Roles = []entities.Role{
				{UserID: newUser.ID, Value: "user"},
			}
			cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
			if err != nil {
				t.Fatal("expected no error", err)
			}
			validToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now(), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
			if err != nil {
				t.Fatal("expected no error", err)
			}
			err = authService.Authorize(validToken, []string{"any"})
			if err != nil {
				t.Fatal("expected no error for 'any' role", err)
			}
		})
	})

	t.Run("authorized role returns no error", func(t *testing.T) {
		authService, _, _, _ := prepare(t)
		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		newUser.Roles = []entities.Role{
			{UserID: newUser.ID, Value: "payments"},
		}
		cc, err := entities.NewCustomClaimsFromValues(newUser.ID, newUser.EarlyAdopter, newUser.Roles, newUser.MetadataPublic)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		validToken, _, err := jwtgen.Generate(cc.ToMap(), time.Now(), baseConfig.JWT.AccessTokenExpirationMin, baseConfig.App.Name, baseConfig.JWT.Secret)
		if err != nil {
			t.Fatal("expected no error", err)
		}
		err = authService.Authorize(validToken, []string{"admin", "user", "payments"})
		if err != nil {
			t.Fatal("expected no error for authorized role", err)
		}
		err = authService.Authorize(validToken, []string{"user"})
		if err.Error() != apperrors.ErrUnauthorizedRole.Error() {
			t.Fatal("expected error ErrUnauthorizedRole", err)
		}
	})
}
