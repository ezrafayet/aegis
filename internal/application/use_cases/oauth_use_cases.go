package usecases

import (
	"fmt"
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/primary"
	"othnx/internal/domain/ports/secondary"
	"othnx/pkg/apperrors"
	"othnx/pkg/tokengen"
	"time"
)

// todo: make this package a generic package for all oauth providers and a factory

type OAuthGithubUseCases struct {
	Config                 entities.Config
	Provider               secondary.OAuthProviderRequests
	UserRepository         secondary.UserRepository
	RefreshTokenRepository secondary.RefreshTokenRepository
	StateRepository        secondary.StateRepository
}

var _ primary.OAuthUseCasesExecutor = OAuthGithubUseCases{}

func NewOAuthGithubUseCases(c entities.Config, p secondary.OAuthProviderRequests, userRepository secondary.UserRepository, refreshTokenRepository secondary.RefreshTokenRepository, stateRepository secondary.StateRepository) OAuthGithubUseCases {
	return OAuthGithubUseCases{
		Config:                 c,
		Provider:               p,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		StateRepository:        stateRepository,
	}
}

// /!\ Can be abused to generate a lot of states - todo: fix
func (s OAuthGithubUseCases) GetAuthURL(redirectUri string) (string, error) {
	state, err := tokengen.Generate("state_", 13)
	if err != nil {
		return "", err
	}
	if err := s.StateRepository.CreateState(entities.NewState(state)); err != nil {
		return "", err
	}
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s",
		s.Config.Auth.Providers.GitHub.ClientID,
		state,
	)
	return redirectURL, nil
}

func (s OAuthGithubUseCases) ExchangeCode(code, state string) (*entities.TokenPair, error) {
	serverState, err := s.StateRepository.GetAndDeleteState(state)
	if err != nil {
		return nil, apperrors.ErrInvalidState
	}
	if serverState.IsExpired() {
		return nil, apperrors.ErrInvalidState
	}
	userInfos, err := s.Provider.GetUserInfos(code, state, "")
	if err != nil {
		return nil, err
	}

	user, err := GetOrCreateUserIfAllowed(s.UserRepository, userInfos, s.Config)
	if err != nil {
		return nil, err
	}

	if user.AuthMethod != "github" {
		return nil, apperrors.ErrWrongAuthMethod
	}

	// todo device-id: pass one, since one session per device is allowed
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := GenerateTokensForUser(user, "device-id", s.Config, s.RefreshTokenRepository)
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

func (s OAuthGithubUseCases) CheckAuthEnabled(provider string) bool {
	switch provider {
	case "github":
		return s.Config.Auth.Providers.GitHub.Enabled
	default:
		return false
	}
}
