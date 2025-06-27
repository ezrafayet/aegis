package usecases

import (
	"fmt"
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/primary_ports"
	"othnx/internal/domain/ports/secondary_ports"
	"othnx/internal/infrastructure/config"
	"othnx/pkg/apperrors"
	"othnx/pkg/tokengen"
	"time"
)

// todo: make this package a generic package for all oauth providers and a factory

type OAuthGithubService struct {
	Config                 config.Config
	Provider               secondaryports.OAuthProviderRequests
	UserRepository         secondaryports.UserRepository
	RefreshTokenRepository secondaryports.RefreshTokenRepository
	StateRepository        secondaryports.StateRepository
}

var _ primaryports.OAuthUseCasesInterface = OAuthGithubService{}

func NewOAuthGithubUseCases(c config.Config, p secondaryports.OAuthProviderRequests, userRepository secondaryports.UserRepository, refreshTokenRepository secondaryports.RefreshTokenRepository, stateRepository secondaryports.StateRepository) OAuthGithubService {
	return OAuthGithubService{
		Config:                 c,
		Provider:               p,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		StateRepository:        stateRepository,
	}
}

// /!\ Can be abused to generate a lot of states - todo: fix
func (s OAuthGithubService) GetAuthURL(redirectUri string) (string, error) {
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

func (s OAuthGithubService) ExchangeCode(code, state string) (*entities.TokenPair, error) {
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

func (s OAuthGithubService) CheckAuthEnabled(provider string) bool {
	switch provider {
	case "github":
		return s.Config.Auth.Providers.GitHub.Enabled
	default:
		return false
	}
}
