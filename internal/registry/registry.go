package registry

import (
	"aegix/internal/domain"
	"aegix/internal/providers/github"
)

type Registry struct {
	GitHubRouter github.OAuthGithubRouter
}

func NewRegistry(c domain.Config) Registry {
	githubProvider := github.NewOAuthGithubProvider(c)
	githubServices := github.NewOAuthGithubService(c, githubProvider)
	githubHandlers := github.NewOAuthGithubHandlers(c, githubServices)
	githubMiddlewares := github.NewOAuthGithubMiddlewares(c)
	githubRouter := github.NewOAuthGithubRouter(githubHandlers, githubMiddlewares)

	return Registry{
		GitHubRouter: githubRouter,
	}
}
