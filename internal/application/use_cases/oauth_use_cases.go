package usecases

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/primary"
	"aegis/internal/domain/ports/secondary"
	"aegis/internal/domain/services"
	"aegis/pkg/apperrors"
	"aegis/pkg/tokengen"
	"fmt"
	"time"
)

// todo: make this package a generic package for all oauth providers and a factory

type OAuthUseCases struct {
	Config                 entities.Config
	Provider               secondary.OAuthProviderRequests
	UserRepository         secondary.UserRepository
	RefreshTokenRepository secondary.RefreshTokenRepository
	StateRepository        secondary.StateRepository
	UserService            *services.UserService
	TokenService           *services.TokenService
}

var _ primary.OAuthUseCasesExecutor = OAuthUseCases{}

func NewOAuthGithubUseCases(c entities.Config, p secondary.OAuthProviderRequests, userRepository secondary.UserRepository, refreshTokenRepository secondary.RefreshTokenRepository, stateRepository secondary.StateRepository) OAuthUseCases {
	userService := services.NewUserService(userRepository, c)
	tokenService := services.NewTokenService(refreshTokenRepository, c)
	return OAuthUseCases{
		Config:                 c,
		Provider:               p,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		StateRepository:        stateRepository,
		UserService:            userService,
		TokenService:           tokenService,
	}
}

// /!\ Can be abused to generate a lot of states - todo: fix
func (s OAuthUseCases) GetAuthURL(redirectUri string) (string, error) {
	state, err := tokengen.Generate("state_", 13)
	if err != nil {
		return "", err
	}
	if err := s.StateRepository.CreateState(entities.NewState(state)); err != nil {
		return "", err
	}
	// todo: to support multiple providers, we need to use a map of providers and their config
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s",
		s.Config.Auth.Providers.GitHub.ClientID,
		state,
	)
	return redirectURL, nil
}

func (s OAuthUseCases) ExchangeCode(code, state string) (*entities.TokenPair, error) {
	serverState, err := s.StateRepository.GetAndDeleteState(state)
	if err != nil {
		return nil, apperrors.ErrInvalidState
	}
	if serverState.IsExpired() {
		return nil, apperrors.ErrInvalidState
	}
	userInfos, err := s.Provider.ExchangeCodeForUserInfos(code, state, "")
	if err != nil {
		return nil, err
	}

	user, err := s.UserService.GetOrCreateUserIfAllowed(userInfos, s.Provider.GetName())
	if err != nil {
		return nil, err
	}

	if user.AuthMethod != s.Provider.GetName() {
		return nil, apperrors.ErrWrongAuthMethod
	}

	// todo device-id: pass one, since one session per device is allowed
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := s.TokenService.GenerateTokensForUser(user, "device-id")
	if err != nil {
		return nil, err
	}

	return &entities.TokenPair{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  time.Unix(atExpiresAt, 0),
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: time.Unix(rtExpiresAt, 0),
	}, nil
}

func (s OAuthUseCases) CheckAuthEnabled() bool {
	// todo: move this logic somewhere more relevant
	switch s.Provider.GetName() {
	case "github":
		return s.Config.Auth.Providers.GitHub.Enabled
	default:
		return false
	}
}
