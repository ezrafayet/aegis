package primary

import "othnx/internal/domain/entities"

type UseCasesExecutor interface {
	// For handlers
	GetSession(accessToken string) (entities.Session, error)
	Logout(refreshToken string) (*entities.TokenPair, error)

	// For middlewares
	CheckAndRefreshToken(accessToken, refreshToken string, forceRefresh bool) (*entities.TokenPair, error)
}
