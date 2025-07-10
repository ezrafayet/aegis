package primary

import "aegis/internal/domain/entities"

type OAuthUseCasesForHandlers interface {
	GetAuthURL(redirectUri string) (string, error)
	ExchangeCode(code, state string) (*entities.TokenPair, error)
}

type OAuthUseCasesForMiddlewares interface {
	CheckAuthEnabled() bool
}

type OAuthUseCasesInterface interface {
	OAuthUseCasesForHandlers
	OAuthUseCasesForMiddlewares
}
