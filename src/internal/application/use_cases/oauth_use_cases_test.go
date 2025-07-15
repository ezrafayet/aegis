package usecases

// import (
// 	"aegis/internal/domain/entities"
// 	"aegis/internal/domain/ports/secondary_ports"
// 	"aegis/internal/infrastructure/config"
// 	"aegis/internal/infrastructure/repositories"
// 	"aegis/pkg/apperrors"
// 	"testing"
// 	"time"

// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// type MockRepository struct{}

// func (p MockRepository) GetUserInfos(code, state, redirectUri string) (*entities.UserInfos, error) {
// 	userInfos := entities.UserInfos{
// 		Name:   "some-name",
// 		Email:  "some-email",
// 		Avatar: "some-avatar",
// 	}
// 	return &userInfos, nil
// }

// func TestOAuthGithubService_ExchangeCode(t *testing.T) {
// 	baseConfig := config.Config{
// 		JWT: config.JWTConfig{
// 			Secret:                     "some-secret",
// 			AccessTokenExpirationMin:   1,
// 			RefreshTokenExpirationDays: 1,
// 		},
// 	}
// 	prepare := func(t *testing.T) (OAuthGithubService, secondaryports.UserRepository, secondaryports.RefreshTokenRepository, *gorm.DB) {
// 		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		db.AutoMigrate(&entities.User{}, &entities.RefreshToken{}, &entities.State{}, &entities.Role{})
// 		refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
// 		oauthProviderRepository := MockRepository{}
// 		userRepository := repositories.NewUserRepository(db)
// 		stateRepository := repositories.NewStateRepository(db)
// 		authService := NewOAuthGithubService(baseConfig, oauthProviderRepository, &userRepository, &refreshTokenRepository, &stateRepository)
// 		return authService, &userRepository, &refreshTokenRepository, db
// 	}
// 	t.Run("should create a user if no user exists", func(t *testing.T) {
// 		ghService, _, _, db := prepare(t)
// 		db.Create(&entities.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
// 		_, _, err := ghService.ExchangeCode("some-code", "some-state")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		var user entities.User
// 		db.First(&user, "email = ?", "some-email")
// 		if user.Email != "some-email" {
// 			t.Fatal("expected user to be created", user.Email)
// 		}
// 		var roles []entities.Role
// 		db.Model(&entities.Role{}).Where("user_id = ?", user.ID).Find(&roles)
// 		if len(roles) != 1 {
// 			t.Fatal("expected user to have 1 role", len(roles))
// 		}
// 		if roles[0].Value != "user" {
// 			t.Fatal("expected user to have role user", roles[0].Value)
// 		}
// 	})
// 	t.Run("should not create a user if user already exists", func(t *testing.T) {
// 		ghService, _, _, db := prepare(t)
// 		db.Create(&entities.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
// 		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		db.Create(&newUser)
// 		_, _, err = ghService.ExchangeCode("some-code", "some-state")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		var count int64
// 		db.Model(&entities.User{}).Where("email = ?", "some-email").Count(&count)
// 		if count != 1 {
// 			t.Fatal("expected user to be created", count)
// 		}
// 	})
// 	t.Run("should not return tokens if user is blocked", func(t *testing.T) {
// 		ghService, _, _, db := prepare(t)
// 		db.Create(&entities.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
// 		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		now := time.Now()
// 		newUser.BlockedAt = &now
// 		db.Create(&newUser)
// 		_, _, err = ghService.ExchangeCode("some-code", "some-state")
// 		if err.Error() != apperrors.ErrUserBlocked.Error() {
// 			t.Fatal("expected error ErrUserBlocked", err)
// 		}
// 	})
// 	t.Run("should not return tokens if user is deleted", func(t *testing.T) {
// 		ghService, _, _, db := prepare(t)
// 		db.Create(&entities.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
// 		newUser, err := entities.NewUser("some-name", "some-avatar", "some-email", "github")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		now := time.Now()
// 		newUser.DeletedAt = &now
// 		db.Create(&newUser)
// 		_, _, err = ghService.ExchangeCode("some-code", "some-state")
// 		if err.Error() != apperrors.ErrUserDeleted.Error() {
// 			t.Fatal("expected error ErrUserDeleted", err)
// 		}
// 	})
// 	t.Run("should issue tokens if user does not have a valid refresh token", func(t *testing.T) {
// 		ghService, _, _, db := prepare(t)
// 		db.Create(&entities.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
// 		at, rt, err := ghService.ExchangeCode("some-code", "some-state")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		var user entities.User
// 		db.First(&user, "email = ?", "some-email")
// 		if user.Email != "some-email" {
// 			t.Fatal("expected user to be created", user.Email)
// 		}
// 		if at == nil || at.Value == "" {
// 			t.Fatal("expected access token to be nil", at)
// 		}
// 		if rt == nil || rt.Value == "" {
// 			t.Fatal("expected refresh token to be nil", rt)
// 		}
// 		var refreshToken entities.RefreshToken
// 		db.First(&refreshToken, "user_id = ?", user.ID)
// 		if refreshToken.Token == "" {
// 			t.Fatal("expected refresh token to be created", refreshToken.Token)
// 		}
// 	})
// 	t.Run("should issue new tokens if user already has a valid refresh token, and delete the previous one", func(t *testing.T) {
// 		ghService, _, _, db := prepare(t)
// 		db.Create(&entities.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
// 		_, rt1, err := ghService.ExchangeCode("some-code", "some-state")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		var user entities.User
// 		db.First(&user, "email = ?", "some-email")
// 		if user.Email != "some-email" {
// 			t.Fatal("expected user to be created", user.Email)
// 		}
// 		if rt1 == nil || rt1.Value == "" {
// 			t.Fatal("expected refresh token to be nil", rt1)
// 		}
// 		db.Create(&entities.State{Value: "some-state", ExpiresAt: time.Now().Add(3 * time.Minute)})
// 		_, rt2, err := ghService.ExchangeCode("some-code", "some-state")
// 		if err != nil {
// 			t.Fatal("expected no error", err)
// 		}
// 		if rt2 == nil || rt2.Value == "" {
// 			t.Fatal("expected refresh token to be nil", rt2)
// 		}
// 		var count int64
// 		db.Model(&entities.RefreshToken{}).Where("user_id = ?", user.ID).Count(&count)
// 		if count != 1 {
// 			t.Fatal("expected refresh token to be deleted", count)
// 		}
// 		var refreshToken entities.RefreshToken
// 		db.First(&refreshToken, "user_id = ?", user.ID)
// 		if refreshToken.Token == rt1.Value {
// 			t.Fatal("expected refresh token to be deleted", refreshToken.Token)
// 		}
// 		if refreshToken.Token != rt2.Value {
// 			t.Fatal("expected refresh token to be created", refreshToken.Token)
// 		}
// 	})
// 	// todo: add specific tests when device id is handled
// }
