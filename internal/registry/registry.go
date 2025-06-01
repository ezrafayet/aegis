package registry

import (
	"aegix/internal/domain"
	"aegix/internal/providers/github"
	"aegix/internal/repository/pg"

	"gorm.io/gorm"
)

type Registry struct {
	GitHubRouter github.OAuthGithubRouter
}

func NewRegistry(c domain.Config, db *gorm.DB) Registry {
	userRepository := pgrepository.NewUserRepository(db)
	refreshTokenRepository := pgrepository.NewRefreshTokenRepository(db)

	githubProvider := github.NewOAuthGithubProvider(c)
	githubServices := github.NewOAuthGithubService(c, githubProvider, &userRepository, &refreshTokenRepository)
	githubHandlers := github.NewOAuthGithubHandlers(c, githubServices)
	githubMiddlewares := github.NewOAuthGithubMiddlewares(c)
	githubRouter := github.NewOAuthGithubRouter(githubHandlers, githubMiddlewares)

	return Registry{
		GitHubRouter: githubRouter,
	}
}
