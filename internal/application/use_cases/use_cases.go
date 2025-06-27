package usecases

import (
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/primary_ports"
	"othnx/internal/domain/ports/secondary_ports"
	"othnx/internal/infrastructure/config"
	"othnx/pkg/apperrors"
	"othnx/pkg/jwtgen"
	"time"
)

type AuthService struct {
	Config                 config.Config
	RefreshTokenRepository secondaryports.RefreshTokenRepository
	UserRepository         secondaryports.UserRepository
}

var _ primaryports.UseCasesInterface = &AuthService{}

func NewService(c config.Config, r secondaryports.RefreshTokenRepository, u secondaryports.UserRepository) AuthService {
	return AuthService{
		Config:                 c,
		RefreshTokenRepository: r,
		UserRepository:         u,
	}
}

func (s AuthService) GetSession(accessToken string) (entities.Session, error) {
	customClaims, err := jwtgen.ReadClaims(accessToken, s.Config)
	if err != nil {
		return entities.Session{}, err
	}
	return entities.Session{
		CustomClaims: customClaims,
	}, nil
}

func (s AuthService) eraseTokens(err error) (*entities.TokenPair, error) {
	return &entities.TokenPair{
		AccessToken:           "",
		AccessTokenExpiresAt:  time.Now(),
		RefreshToken:          "",
		RefreshTokenExpiresAt: time.Now(),
	}, err
}

func (s AuthService) Logout(refreshToken string) (*entities.TokenPair, error) {
	if refreshToken != "" {
		_ = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	}
	return s.eraseTokens(nil)
}

func (s AuthService) CheckAndRefreshToken(accessToken, refreshToken string, forceRefresh bool) (*entities.TokenPair, error) {
	_, err := jwtgen.ReadClaims(accessToken, s.Config)
	if err == nil && !forceRefresh {
		return s.eraseTokens(nil)
	}
	if err != nil && err.Error() != apperrors.ErrAccessTokenExpired.Error() {
		return s.eraseTokens(err)
	}
	refreshTokenObject, err := s.RefreshTokenRepository.GetRefreshTokenByToken(refreshToken)
	if err != nil {
		return s.eraseTokens(err)
	}
	if refreshTokenObject.IsExpired() {
		return s.eraseTokens(apperrors.ErrRefreshTokenExpired)
	}
	// todo: check device id
	user, err := s.UserRepository.GetUserByID(refreshTokenObject.UserID)
	if err != nil {
		return s.eraseTokens(err)
	}
	if user.IsDeleted() {
		return s.eraseTokens(apperrors.ErrUserDeleted)
	}
	if user.IsBlocked() {
		return s.eraseTokens(apperrors.ErrUserBlocked)
	}
	if s.Config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return s.eraseTokens(apperrors.ErrEarlyAdoptersOnly)
	}

	err = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		return s.eraseTokens(err)
	}
	// todo device-id: pass one, since one session per device is allowed
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := GenerateTokensForUser(user, "device-id", s.Config, s.RefreshTokenRepository)
	if err != nil {
		return s.eraseTokens(err)
	}
	return &entities.TokenPair{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  time.Unix(atExpiresAt, 0),
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: time.Unix(rtExpiresAt, 0),
	}, nil
}
