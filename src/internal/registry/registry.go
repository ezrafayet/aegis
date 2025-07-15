package registry

import (
	usecases "aegis/internal/application/use_cases"
	"aegis/internal/domain/entities"
	"aegis/internal/infrastructure/handlers"
	"aegis/internal/infrastructure/middlewares"
	"aegis/internal/infrastructure/repositories"
	"aegis/pkg/plugins/providers/discord"
	"aegis/pkg/plugins/providers/github"
	"fmt"

	"gorm.io/gorm"
)

type Registry struct {
	Handlers    handlers.HandlersInterface
	Middlewares middlewares.AuthMiddlewareInterface
	Providers   []Provider
}

func NewRegistry(c entities.Config, db *gorm.DB) (Registry, error) {
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	stateRepository := repositories.NewStateRepository(db)

	authService := usecases.NewService(c, refreshTokenRepository, userRepository)
	authHandlers := handlers.NewHandlers(c, authService)
	authMiddlewares := middlewares.NewAuthMiddleware(c, authService)

	providers := []Provider{
		NewProvider(
			c, github.NewOAuthGithubRepository(
				c.Auth.Providers.GitHub.Enabled,
				c.Auth.Providers.GitHub.ClientID,
				c.Auth.Providers.GitHub.ClientSecret,
				fmt.Sprintf("%s/auth/github/callback", c.App.URL)),
			userRepository,
			refreshTokenRepository,
			stateRepository),
		NewProvider(
			c, discord.NewOAuthDiscordRepository(
				c.Auth.Providers.Discord.Enabled,
				c.Auth.Providers.Discord.ClientID,
				c.Auth.Providers.Discord.ClientSecret,
				fmt.Sprintf("%s/auth/discord/callback", c.App.URL)),
			userRepository,
			refreshTokenRepository,
			stateRepository),
	}

	return Registry{
		Handlers:    authHandlers,
		Middlewares: authMiddlewares,
		Providers:   providers,
	}, nil
}
