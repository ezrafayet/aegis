package registry

import (
	"othnx/internal/infrastructure/config"
	"othnx/internal/infrastructure/handlers"
	"othnx/internal/infrastructure/middlewares"
	"othnx/internal/infrastructure/repositories"
	"othnx/internal/application/use_cases"
	"othnx/internal/infrastructure/providers/github"

	"gorm.io/gorm"
)

type Registry struct {
	Handlers handlers.HandlersInterface
	Middlewares middlewares.AuthMiddlewareInterface
	OAuthHandlers handlers.OAuthGithubHandlers
	OAuthMiddlewares middlewares.OAuthGithubMiddlewares
}

func NewRegistry(c config.Config, db *gorm.DB) Registry {
	userRepository := repositories.NewUserRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(db)
	stateRepository := repositories.NewStateRepository(db)

	usecases := usecases.NewService(c, &refreshTokenRepository, &userRepository)
	handlers := handlers.NewHandlers(c, usecases)
	middlewares := middlewares.NewAuthMiddleware(c, usecases)

	githubProvider := github.NewOAuthGithubRepository(c)
	githubUsecases := usecases.NewOAuthGithubUseCases(c, githubProvider, &userRepository, &refreshTokenRepository, &stateRepository)
	githubHandlers := handlers.NewOAuthGithubHandlers(c, githubUsecases)
	githubMiddlewares := middlewares.NewOAuthGithubMiddlewares(c)

	return Registry{
		Handlers: handlers,
		Middlewares: middlewares,
		OAuthHandlers: githubHandlers,
		OAuthMiddlewares: githubMiddlewares,
	}
}
