package usecases

import (
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/primary"
	"othnx/internal/domain/ports/secondary"
	"othnx/internal/domain/services"
	"othnx/pkg/apperrors"
	"othnx/pkg/jwtgen"
	"time"
)

type UseCases struct {
	Config                 entities.Config
	RefreshTokenRepository secondary.RefreshTokenRepository
	UserRepository         secondary.UserRepository
	TokenService           *services.TokenService
}

var _ primary.UseCasesExecutor = &UseCases{}

func NewService(c entities.Config, r secondary.RefreshTokenRepository, u secondary.UserRepository) UseCases {
	tokenService := services.NewTokenService(r, c)
	return UseCases{
		Config:                 c,
		RefreshTokenRepository: r,
		UserRepository:         u,
		TokenService:           tokenService,
	}
}

func (s UseCases) GetSession(accessToken string) (entities.Session, error) {
	customClaims, err := jwtgen.ReadClaims(accessToken, s.Config)
	if err != nil {
		return entities.Session{}, err
	}
	return entities.Session{
		CustomClaims: customClaims,
	}, nil
}

func (s UseCases) eraseTokens(err error) (*entities.TokenPair, error) {
	return &entities.TokenPair{
		AccessToken:           "",
		AccessTokenExpiresAt:  time.Now(),
		RefreshToken:          "",
		RefreshTokenExpiresAt: time.Now(),
	}, err
}

func (s UseCases) Logout(refreshToken string) (*entities.TokenPair, error) {
	if refreshToken != "" {
		_ = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	}
	return s.eraseTokens(nil)
}

func (s UseCases) CheckAndRefreshToken(accessToken, refreshToken string, forceRefresh bool) (*entities.TokenPair, error) {
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
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := s.TokenService.GenerateTokensForUser(user, "device-id")
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
