package auth

import (
	"aegix/internal/domain"
	"aegix/internal/repository"
)

type AuthServiceInterface interface {
	GetSession(accessToken string) (string, error)
	Logout(refreshToken string) error
}

type AuthService struct {
	Config domain.Config
	RefreshTokenRepository repository.RefreshTokenRepository
	UserRepository repository.UserRepository
}

var _ AuthServiceInterface = &AuthService{}

func NewAuthService(c domain.Config, r repository.RefreshTokenRepository, u repository.UserRepository) AuthService {
	return AuthService{
		Config: c,
		RefreshTokenRepository: r,
		UserRepository: u,
	}
}

func (s AuthService) GetSession(accessToken string) (string, error) {
	return "", nil
}

func (s AuthService) Logout(refreshToken string) error {
	return nil
}
