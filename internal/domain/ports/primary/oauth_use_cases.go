package primary

import "aegis/internal/domain/entities"

type OAuthUseCasesExecutor interface {
	// For handlers
	GetAuthURL(redirectUri string) (string, error)
	ExchangeCode(code, state string) (*entities.TokenPair, error)

	// For middlewares
	CheckAuthEnabled(provider string) bool
}
