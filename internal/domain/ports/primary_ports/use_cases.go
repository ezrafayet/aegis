package primary_ports

import "othnx/internal/domain/entities"

type UseCasesInterface interface {
	// For handlers
	GetSession(accessToken string) (entities.Session, error)
	Logout(refreshToken string) (*entities.TokenPair, error)

	// For middlewares
	CheckAndRefreshToken(accessToken, refreshToken string, forceRefresh bool) (*entities.TokenPair, error)
}
