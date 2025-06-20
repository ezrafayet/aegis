package registry

import (
	"othnx/internal/components/auth"
	"othnx/internal/components/providers/github"
	"othnx/internal/domain"
	"othnx/internal/repository"

	"gorm.io/gorm"
)

type Registry struct {
	GitHubRouter github.OAuthGithubRouter
	AuthRouter   auth.AuthRouter
}

func NewRegistry(c domain.Config, db *gorm.DB) Registry {
	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	stateRepository := repository.NewStateRepository(db)

	authService := auth.NewAuthService(c, refreshTokenRepository, userRepository)
	authHandlers := auth.NewAuthHandlers(c, authService)
	authMiddlewares := auth.NewAuthMiddleware(c, authService)
	authRouter := auth.NewAuthRouter(authHandlers, authMiddlewares)

	githubProvider := github.NewOAuthGithubRepository(c)
	githubServices := github.NewOAuthGithubService(c, githubProvider, &userRepository, &refreshTokenRepository, &stateRepository)
	githubHandlers := github.NewOAuthGithubHandlers(c, githubServices)
	githubMiddlewares := github.NewOAuthGithubMiddlewares(c)
	githubRouter := github.NewOAuthGithubRouter(githubHandlers, githubMiddlewares)

	return Registry{
		GitHubRouter: githubRouter,
		AuthRouter:   authRouter,
	}
}
