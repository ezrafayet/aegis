package registry

import (
	"aegix/internal/components/auth"
	"aegix/internal/components/providers/github"
	"aegix/internal/domain"
	"aegix/internal/repository"

	"gorm.io/gorm"
)

type Registry struct {
	GitHubRouter github.OAuthGithubRouter
	AuthRouter   auth.AuthRouter
}

func NewRegistry(c domain.Config, db *gorm.DB) Registry {
	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)

	authService := auth.NewAuthService(c, refreshTokenRepository, userRepository)
	authHandlers := auth.NewAuthHandlers(c, authService)
	authMiddlewares := auth.NewAuthMiddleware(c, authService)
	authRouter := auth.NewAuthRouter(authHandlers, authMiddlewares)

	githubProvider := github.NewOAuthGithubProvider(c)
	githubServices := github.NewOAuthGithubService(c, githubProvider, &userRepository, &refreshTokenRepository)
	githubHandlers := github.NewOAuthGithubHandlers(c, githubServices)
	githubMiddlewares := github.NewOAuthGithubMiddlewares(c)
	githubRouter := github.NewOAuthGithubRouter(githubHandlers, githubMiddlewares)

	return Registry{
		GitHubRouter: githubRouter,
		AuthRouter:   authRouter,
	}
}
