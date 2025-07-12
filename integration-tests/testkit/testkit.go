package testkit

import (
	"aegis/internal/domain/entities"
	"aegis/internal/infrastructure/database"
	"aegis/internal/registry"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgresContainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuite struct {
	T           *testing.T
	CTX         context.Context
	PgContainer *postgresContainer.PostgresContainer
	Db          *gorm.DB
	Server      *httptest.Server
	MockGithub  *httptest.Server
	MockDiscord *httptest.Server
	Config      entities.Config
}

func SetupTestSuite(t *testing.T) *TestSuite {
	ctx := context.Background()
	pgContainer, db, connStr := setupDatabase(t, ctx)
	suite := &TestSuite{
		T:           t,
		CTX:         ctx,
		PgContainer: pgContainer,
		Db:          db,
	}
	suite.setupMockServers()
	suite.setupConfig(connStr)
	suite.setupTestServer()
	return suite
}

func setupDatabase(t *testing.T, ctx context.Context) (*postgresContainer.PostgresContainer, *gorm.DB, string) {
	pgContainer, err := postgresContainer.Run(ctx,
		"postgres:15-alpine",
		postgresContainer.WithDatabase("test_db"),
		postgresContainer.WithUsername("test_user"),
		postgresContainer.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	require.NoError(t, err)
	err = database.Migrate(db)
	require.NoError(t, err)
	return pgContainer, db, connStr
}

func (s *TestSuite) setupMockServers() {
	// Mock GitHub OAuth server
	s.MockGithub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login/oauth/access_token":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"access_token": "mock_access_token",
				"token_type":   "bearer",
				"scope":        "user:email",
			})
		case "/user":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"login":      "testuser",
				"name":       "Test User",
				"email":      "test@example.com",
				"avatar_url": "https://avatars.githubusercontent.com/u/123?v=4",
			})
		case "/user/emails":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"email":    "test@example.com",
					"primary":  true,
					"verified": true,
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	// Mock Discord OAuth server
	s.MockDiscord = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/oauth2/token":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"access_token": "mock_discord_access_token",
				"token_type":   "Bearer",
				"scope":        "identify email",
			})
		case "/api/users/@me":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":            "123456789",
				"username":      "testuser",
				"email":         "test@example.com",
				"avatar":        "a_1234567890abcdef",
				"discriminator": "0001",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func (s *TestSuite) setupConfig(dbURL string) {
	s.Config = entities.Config{}

	// App configuration
	s.Config.App.Name = "TestApp"
	s.Config.App.URL = "http://localhost:8080"
	s.Config.App.CorsAllowedOrigins = []string{"http://localhost:8080"}
	s.Config.App.EarlyAdoptersOnly = false
	s.Config.App.RedirectAfterSuccess = "http://localhost:8080/login-success"
	s.Config.App.RedirectAfterError = "http://localhost:8080/login-error"
	s.Config.App.InternalAPIKeys = []string{"test-api-key"}
	s.Config.App.Port = 8080

	// Database configuration
	s.Config.DB.PostgresURL = dbURL

	// JWT configuration
	s.Config.JWT.Secret = "test-jwt-secret-key-for-testing"
	s.Config.JWT.AccessTokenExpirationMin = 15
	s.Config.JWT.RefreshTokenExpirationDays = 7

	// Auth providers configuration
	s.Config.Auth.Providers.GitHub.Enabled = true
	s.Config.Auth.Providers.GitHub.AppName = "TestApp"
	s.Config.Auth.Providers.GitHub.ClientID = "test-github-client-id"
	s.Config.Auth.Providers.GitHub.ClientSecret = "test-github-client-secret"

	s.Config.Auth.Providers.Discord.Enabled = true
	s.Config.Auth.Providers.Discord.AppName = "TestApp"
	s.Config.Auth.Providers.Discord.ClientID = "test-discord-client-id"
	s.Config.Auth.Providers.Discord.ClientSecret = "test-discord-client-secret"

	// Cookies configuration
	s.Config.Cookies.Domain = "localhost"
	s.Config.Cookies.Secure = false
	s.Config.Cookies.HTTPOnly = true
	s.Config.Cookies.SameSite = 1
	s.Config.Cookies.Path = "/"

	// User configuration
	s.Config.User.Roles = []string{"user", "platform_admin"}
}

func (s *TestSuite) setupTestServer() {
	r, err := registry.NewRegistry(s.Config, s.Db)
	require.NoError(s.T, err)
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	group := e.Group("/auth")
	group.GET("/me", r.Handlers.GetSession, r.Middlewares.CheckAndRefreshToken)
	group.GET("/refresh", r.Handlers.DoNothing, r.Middlewares.CheckAndForceRefreshToken)
	group.GET("/logout", r.Handlers.Logout)
	group.GET("/health", r.Handlers.DoNothing)
	for _, provider := range r.Providers {
		group.GET(fmt.Sprintf("/%s", provider.Name), provider.Handlers.GetAuthURL, provider.Middlewares.CheckAuthEnabled)
		group.GET(fmt.Sprintf("/%s/callback", provider.Name), provider.Handlers.ExchangeCode, provider.Middlewares.CheckAuthEnabled)
	}
	s.Server = httptest.NewServer(e)
}

func (s *TestSuite) Teardown() {
	if s.Server != nil {
		s.Server.Close()
	}
	if s.MockGithub != nil {
		s.MockGithub.Close()
	}
	if s.MockDiscord != nil {
		s.MockDiscord.Close()
	}
	if s.PgContainer != nil {
		s.PgContainer.Terminate(s.CTX)
	}
}

func (s *TestSuite) CreateUser(t *testing.T, user entities.User, roles []string) entities.User {
	err := s.Db.Model(&entities.User{}).Create(&user).Error
	require.NoError(t, err)
	userRoles := []entities.Role{}
	if len(roles) > 0 {
		for _, role := range roles {
			err := s.Db.Model(&entities.Role{}).Create(&entities.Role{
				UserID: user.ID,
				Value:  role,
			}).Error
			require.NoError(t, err)
			userRoles = append(userRoles, entities.Role{
				UserID: user.ID,
				Value:  role,
			})
		}
	}
	user.Roles = userRoles
	return user
}

func (s *TestSuite) CreateRefreshToken(t *testing.T, refreshToken entities.RefreshToken) entities.RefreshToken {
	err := s.Db.Model(&entities.RefreshToken{}).Create(&refreshToken).Error
	require.NoError(t, err)
	return refreshToken
}
