package registry

import (
	"aegis/internal/application/use_cases"
	"aegis/internal/domain/entities"
	"aegis/internal/infrastructure/handlers"
	"aegis/internal/infrastructure/middlewares"
	"aegis/internal/infrastructure/providers/discord"
	"aegis/internal/infrastructure/providers/github"
	"aegis/internal/infrastructure/repositories"

	"gorm.io/gorm"
)

type Registry struct {
	Handlers    handlers.HandlersInterface
	Middlewares middlewares.AuthMiddlewareInterface
	Providers   []Provider
}

func NewRegistry(c entities.Config, db *gorm.DB) Registry {
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	stateRepository := repositories.NewStateRepository(db)

	authService := usecases.NewService(c, &refreshTokenRepository, &userRepository)
	authHandlers := handlers.NewHandlers(c, authService)
	authMiddlewares := middlewares.NewAuthMiddleware(c, authService)

	providers := []Provider{
		NewProvider(c, github.NewOAuthGithubRepository(c), &userRepository, &refreshTokenRepository, &stateRepository),
		NewProvider(c, discord.NewOAuthDiscordRepository(c), &userRepository, &refreshTokenRepository, &stateRepository),
	}

	return Registry{
		Handlers:    authHandlers,
		Middlewares: authMiddlewares,
		Providers:   providers,
	}
}
