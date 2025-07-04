package registry

import (
	"aegis/internal/application/use_cases"
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"aegis/internal/infrastructure/handlers"
	"aegis/internal/infrastructure/middlewares"
)

type Provider struct {
	Name        string
	Handlers    handlers.OAuthHandlersInterface
	Middlewares middlewares.OAuthMiddlewaresInterface
}

func NewProvider(
	c entities.Config,
	provider secondary.OAuthProviderInterface,
	userRepository secondary.UserRepository,
	refreshTokenRepository secondary.RefreshTokenRepository,
	stateRepository secondary.StateRepository,
) Provider {
	service := usecases.NewOAuthGithubUseCases(c, provider, userRepository, refreshTokenRepository, stateRepository)
	handlers := handlers.NewOAuthHandlers(c, service)
	middlewares := middlewares.NewOAuthMiddlewares(c)

	return Provider{
		Name:        provider.GetName(),
		Handlers:    handlers,
		Middlewares: middlewares,
	}
}
