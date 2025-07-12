package integration

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
	t           *testing.T
	ctx         context.Context
	pgContainer *postgresContainer.PostgresContainer
	db          *gorm.DB
	server      *httptest.Server
	mockGithub  *httptest.Server
	mockDiscord *httptest.Server
	config      entities.Config
}

func setupTestSuite(t *testing.T) *TestSuite {
	ctx := context.Background()
	pgContainer, db, connStr := setupDatabase(t, ctx)
	suite := &TestSuite{
		t:           t,
		ctx:         ctx,
		pgContainer: pgContainer,
		db:          db,
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
	s.mockGithub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	s.mockDiscord = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	s.config = entities.Config{}

	// App configuration
	s.config.App.Name = "TestApp"
	s.config.App.URL = "http://localhost:8080"
	s.config.App.CorsAllowedOrigins = []string{"http://localhost:8080"}
	s.config.App.EarlyAdoptersOnly = false
	s.config.App.RedirectAfterSuccess = "http://localhost:8080/login-success"
	s.config.App.RedirectAfterError = "http://localhost:8080/login-error"
	s.config.App.InternalAPIKeys = []string{"test-api-key"}
	s.config.App.Port = 8080

	// Database configuration
	s.config.DB.PostgresURL = dbURL

	// JWT configuration
	s.config.JWT.Secret = "test-jwt-secret-key-for-testing"
	s.config.JWT.AccessTokenExpirationMin = 15
	s.config.JWT.RefreshTokenExpirationDays = 7

	// Auth providers configuration
	s.config.Auth.Providers.GitHub.Enabled = true
	s.config.Auth.Providers.GitHub.AppName = "TestApp"
	s.config.Auth.Providers.GitHub.ClientID = "test-github-client-id"
	s.config.Auth.Providers.GitHub.ClientSecret = "test-github-client-secret"

	s.config.Auth.Providers.Discord.Enabled = true
	s.config.Auth.Providers.Discord.AppName = "TestApp"
	s.config.Auth.Providers.Discord.ClientID = "test-discord-client-id"
	s.config.Auth.Providers.Discord.ClientSecret = "test-discord-client-secret"

	// Cookies configuration
	s.config.Cookies.Domain = "localhost"
	s.config.Cookies.Secure = false
	s.config.Cookies.HTTPOnly = true
	s.config.Cookies.SameSite = 1
	s.config.Cookies.Path = "/"

	// User configuration
	s.config.User.Roles = []string{"user", "platform_admin"}
}

func (s *TestSuite) setupTestServer() {
	r, err := registry.NewRegistry(s.config, s.db)
	require.NoError(s.t, err)
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
	s.server = httptest.NewServer(e)
}

func (s *TestSuite) teardown() {
	if s.server != nil {
		s.server.Close()
	}
	if s.mockGithub != nil {
		s.mockGithub.Close()
	}
	if s.mockDiscord != nil {
		s.mockDiscord.Close()
	}
	if s.pgContainer != nil {
		s.pgContainer.Terminate(s.ctx)
	}
}

func (s *TestSuite) CreateUser(t *testing.T, user entities.User, roles []string) entities.User {
	err := s.db.Model(&entities.User{}).Create(&user).Error
	require.NoError(t, err)
	userRoles := []entities.Role{}
	if len(roles) > 0 {
		for _, role := range roles {
			err := s.db.Model(&entities.Role{}).Create(&entities.Role{
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
	err := s.db.Model(&entities.RefreshToken{}).Create(&refreshToken).Error
	require.NoError(t, err)
	return refreshToken
}
