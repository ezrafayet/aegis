package registry

import (
	usecases "aegis/internal/application/use_cases"
	"aegis/internal/domain/entities"
	"aegis/internal/infrastructure/handlers"
	"aegis/internal/infrastructure/middlewares"
	"aegis/internal/infrastructure/providers/github"
	"aegis/internal/infrastructure/repositories"

	"gorm.io/gorm"
)

type Registry struct {
	Handlers         handlers.HandlersInterface
	Middlewares      middlewares.AuthMiddlewareInterface
	OAuthHandlers    handlers.OAuthGithubHandlers
	OAuthMiddlewares middlewares.OAuthGithubMiddlewares
}

func NewRegistry(c entities.Config, db *gorm.DB) Registry {
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	stateRepository := repositories.NewStateRepository(db)

	authService := usecases.NewService(c, &refreshTokenRepository, &userRepository)
	authHandlers := handlers.NewHandlers(c, authService)
	authMiddlewares := middlewares.NewAuthMiddleware(c, authService)

	githubProvider := github.NewOAuthGithubRepository(c)
	githubUsecases := usecases.NewOAuthGithubUseCases(c, githubProvider, &userRepository, &refreshTokenRepository, &stateRepository)
	githubHandlers := handlers.NewOAuthGithubHandlers(c, githubUsecases)
	githubMiddlewares := middlewares.NewOAuthGithubMiddlewares(c)

	return Registry{
		Handlers:         authHandlers,
		Middlewares:      authMiddlewares,
		OAuthHandlers:    githubHandlers,
		OAuthMiddlewares: githubMiddlewares,
	}
}
