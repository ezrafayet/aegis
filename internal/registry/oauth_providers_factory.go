package registry

import (
	"auth302/internal/application/use_cases"
	"auth302/internal/domain/entities"
	"auth302/internal/domain/ports/secondary"
	"auth302/internal/infrastructure/handlers"
	"auth302/internal/infrastructure/middlewares"
	"auth302/pkg/plugins/providers"
)

type Provider struct {
	Name        string
	Handlers    handlers.OAuthHandlersInterface
	Middlewares middlewares.OAuthMiddlewaresInterface
}

func NewProvider(
	c entities.Config,
	provider providers.OAuthProviderInterface,
	userRepository secondary.UserRepository,
	refreshTokenRepository secondary.RefreshTokenRepository,
	stateRepository secondary.StateRepository,
) Provider {
	service := usecases.NewOAuthUseCases(c, provider, userRepository, refreshTokenRepository, stateRepository)
	handlers := handlers.NewOAuthHandlers(c, service)
	middlewares := middlewares.NewOAuthMiddlewares(c, service)

	return Provider{
		Name:        provider.GetName(),
		Handlers:    handlers,
		Middlewares: middlewares,
	}
}
