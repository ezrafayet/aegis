package registry

import (
	"othnx/internal/application/use_cases"
	"othnx/internal/domain/entities"
	"othnx/internal/infrastructure/handlers"
	"othnx/internal/infrastructure/middlewares"
	"othnx/internal/infrastructure/providers/github"
	"othnx/internal/infrastructure/repositories"

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
