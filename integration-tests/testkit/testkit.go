package testkit

import (
	usecases "aegis/internal/application/use_cases"
	"aegis/internal/domain/entities"
	"aegis/internal/infrastructure/database"
	"aegis/internal/infrastructure/handlers"
	"aegis/internal/infrastructure/middlewares"
	"aegis/internal/infrastructure/repositories"
	"aegis/internal/registry"
	"aegis/pkg/urlbuilder"
	"context"
	"fmt"
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
	Config      entities.Config
}

func SetupTestSuite(t *testing.T, config entities.Config) *TestSuite {
	ctx := context.Background()
	pgContainer, db, connStr := setupDatabase(t, ctx)
	suite := &TestSuite{
		T:           t,
		CTX:         ctx,
		PgContainer: pgContainer,
		Db:          db,
	}
	suite.setupMockServers()
	suite.setupConfig(connStr, config)
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
	// No mock servers needed with fake providers!
	// The fake providers handle everything internally
}

func (s *TestSuite) setupConfig(dbURL string, config entities.Config) {
	config.DB.PostgresURL = dbURL
	s.Config = config
}

func GetBaseConfig() entities.Config {
	config := entities.Config{}

	// App configuration
	config.App.Name = "TestApp"
	config.App.URL = "http://localhost:8080"
	config.App.CorsAllowedOrigins = []string{"http://localhost:8080"}
	config.App.EarlyAdoptersOnly = false
	config.App.RedirectAfterSuccess = "http://localhost:8080/login-success"
	config.App.RedirectAfterError = "http://localhost:8080/login-error"
	config.App.InternalAPIKeys = []string{"test-api-key"}
	config.App.Port = 8080

	// Database configuration
	config.DB.PostgresURL = "dburl, replaced by testkit"

	// JWT configuration
	config.JWT.Secret = "test-jwt-secret-key-for-testing"
	config.JWT.AccessTokenExpirationMin = 15
	config.JWT.RefreshTokenExpirationDays = 7

	// Auth providers configuration
	config.Auth.Providers.GitHub.Enabled = true
	config.Auth.Providers.GitHub.AppName = "TestApp"
	config.Auth.Providers.GitHub.ClientID = "test-github-client-id"
	config.Auth.Providers.GitHub.ClientSecret = "test-github-client-secret"

	config.Auth.Providers.Discord.Enabled = true
	config.Auth.Providers.Discord.AppName = "TestApp"
	config.Auth.Providers.Discord.ClientID = "test-discord-client-id"
	config.Auth.Providers.Discord.ClientSecret = "test-discord-client-secret"

	// Cookies configuration
	config.Cookies.Domain = "localhost"
	config.Cookies.Secure = false
	config.Cookies.HTTPOnly = true
	config.Cookies.SameSite = 1
	config.Cookies.Path = "/"

	// User configuration
	config.User.Roles = []string{"user", "platform_admin"}

	return config
}

func (s *TestSuite) setupTestServer() {
	// Create a custom registry with mock server URLs
	r, err := s.createTestRegistry()
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

func (s *TestSuite) createTestRegistry() (registry.Registry, error) {
	userRepository := repositories.NewUserRepository(s.Db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(s.Db)
	stateRepository := repositories.NewStateRepository(s.Db)

	authService := usecases.NewService(s.Config, refreshTokenRepository, userRepository)
	authHandlers := handlers.NewHandlers(s.Config, authService)
	authMiddlewares := middlewares.NewAuthMiddleware(s.Config, authService)

	redirectURLBase, err := urlbuilder.Build(s.Config.App.URL, "/auth/%s/callback", map[string]string{})
	if err != nil {
		return registry.Registry{}, err
	}

	// Use fake providers
	providers := []registry.Provider{
		registry.NewProvider(
			s.Config, NewFakeOAuthProvider(
				"github",
				s.Config.Auth.Providers.GitHub.Enabled,
				s.Config.Auth.Providers.GitHub.ClientID,
				s.Config.Auth.Providers.GitHub.ClientSecret,
				fmt.Sprintf(redirectURLBase, "github")),
			userRepository,
			refreshTokenRepository,
			stateRepository),
		registry.NewProvider(
			s.Config, NewFakeOAuthProvider(
				"discord",
				s.Config.Auth.Providers.Discord.Enabled,
				s.Config.Auth.Providers.Discord.ClientID,
				s.Config.Auth.Providers.Discord.ClientSecret,
				fmt.Sprintf(redirectURLBase, "discord")),
			userRepository,
			refreshTokenRepository,
			stateRepository),
	}

	return registry.Registry{
		Handlers:    authHandlers,
		Middlewares: authMiddlewares,
		Providers:   providers,
	}, nil
}

func (s *TestSuite) Teardown() {
	if s.Server != nil {
		s.Server.Close()
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
