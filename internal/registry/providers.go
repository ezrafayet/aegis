package registry

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"aegis/internal/infrastructure/handlers"
	"aegis/internal/infrastructure/middlewares"
	"aegis/internal/application/use_cases"
)

type Provider struct {
	Name string
	Handlers handlers.OAuthHandlersInterface
	Middlewares middlewares.OAuthMiddlewaresInterface
}

func NewProvider(
	c entities.Config,
	requests secondary.OAuthProviderRequests,
	userRepository secondary.UserRepository,
	refreshTokenRepository secondary.RefreshTokenRepository,
	stateRepository secondary.StateRepository,
) Provider {
	service := usecases.NewOAuthGithubUseCases(c, requests, userRepository, refreshTokenRepository, stateRepository)
	handlers := handlers.NewOAuthHandlers(c, service)
	middlewares := middlewares.NewOAuthMiddlewares(c)

	return Provider{
		Name: requests.GetName(),
		Handlers: handlers,
		Middlewares: middlewares,
	}
}
