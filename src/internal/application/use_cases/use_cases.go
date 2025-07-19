package usecases

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/internal/domain/ports/secondary"
	"aegis/internal/domain/services"
	"aegis/pkg/apperrors"
	"aegis/pkg/jwtgen"
	"slices"
	"strings"
	"time"
)

type UseCases struct {
	Config                 entities.Config
	RefreshTokenRepository secondary.RefreshTokenRepository
	UserRepository         secondary.UserRepository
	TokenService           *services.TokenService
}

var _ primary.UseCasesInterface = (*UseCases)(nil)

func NewService(c entities.Config, r secondary.RefreshTokenRepository, u secondary.UserRepository) *UseCases {
	tokenService := services.NewTokenService(r, c)
	return &UseCases{
		Config:                 c,
		RefreshTokenRepository: r,
		UserRepository:         u,
		TokenService:           tokenService,
	}
}

func (s UseCases) GetSession(accessToken string) (entities.Session, error) {
	ccMap, err := jwtgen.ReadClaims(accessToken, s.Config.JWT.Secret)
	if err != nil {
		return entities.Session{}, err
	}
	cc, err := entities.NewCusomClaimsFromMap(ccMap)
	if err != nil {
		return entities.Session{}, err
	}
	return entities.Session{
		CustomClaims: *cc,
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
	if accessToken != "" && !forceRefresh {
		_, err := jwtgen.ReadClaims(accessToken, s.Config.JWT.Secret)
		if err == nil {
			return nil, nil
		}
		if err.Error() != apperrors.ErrAccessTokenInvalid.Error() && err.Error() != apperrors.ErrAccessTokenExpired.Error() {
			return nil, err
		}
	}
	if refreshToken == "" {
		return nil, apperrors.ErrRefreshTokenInvalid
	}
	refreshTokenObject, err := s.RefreshTokenRepository.GetRefreshTokenByToken(refreshToken)
	if err != nil {
		if refreshTokenObject.Token == "" {
			return nil, apperrors.ErrRefreshTokenInvalid
		}
		return nil, err
	}
	if refreshTokenObject.IsExpired() {
		return nil, apperrors.ErrRefreshTokenExpired
	}
	// todo: check device id
	user, err := s.UserRepository.GetUserByID(refreshTokenObject.UserID)
	if err != nil {
		return nil, err
	}
	if user.IsDeleted() {
		return nil, apperrors.ErrUserDeleted
	}
	if user.IsBlocked() {
		return nil, apperrors.ErrUserBlocked
	}
	if s.Config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return nil, apperrors.ErrEarlyAdoptersOnly
	}

	err = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// todo device-id: pass one, since one session per device is allowed
	newAccessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := s.TokenService.GenerateTokensForUser(user, "device-id")
	if err != nil {
		return nil, err
	}
	return &entities.TokenPair{
		AccessToken:           newAccessToken,
		AccessTokenExpiresAt:  time.Unix(atExpiresAt, 0),
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: time.Unix(rtExpiresAt, 0),
	}, nil
}

func (s UseCases) Authorize(accessToken string, authorizedRoles []string) (*entities.CustomClaims, error) {
	if len(authorizedRoles) == 0 {
		return nil, apperrors.ErrNoRoles
	}
	ccMap, err := jwtgen.ReadClaims(accessToken, s.Config.JWT.Secret)
	if err != nil {
		return nil, err
	}
	cc, err := entities.NewCusomClaimsFromMap(ccMap)
	if err != nil {
		return nil, err
	}
	if slices.Contains(authorizedRoles, "any") {
		return cc, nil
	}
	authorized := false
	for _, authorizedRole := range authorizedRoles {
		if ccMap["roles"] != nil && strings.Contains(ccMap["roles"].(string), authorizedRole) {
			authorized = true
			break
		}
	}
	if !authorized {
		return nil, apperrors.ErrUnauthorizedRole
	}
	return cc, nil
}

func (s UseCases) AuthorizeInternalAPICall(key string) error {
	key = strings.TrimPrefix(key, "Bearer ")
	for _, k := range s.Config.App.InternalAPIKeys {
		if k == key {
			return nil
		}
	}
	return apperrors.ErrInternalAPIKeyInvalid
}
