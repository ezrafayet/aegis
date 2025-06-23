package primary_ports

import "othnx/internal/domain/entities"

type OAuthUseCasesInterface interface {
	// For handlers
	GetAuthURL(redirectUri string) (string, error)
	ExchangeCode(code, state string) (*entities.TokenPair, error)

	// For middlewares
	CheckAuthEnabled(provider string) bool
}
