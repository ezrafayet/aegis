package github

import (
	"fmt"
	"net/http"
	"othnx/internal/components/providers/providersports"
	"othnx/internal/domain"
	"othnx/pkg/cookies"
)

type OAuthGithubService struct {
	Config                 domain.Config
	Provider               providersports.OAuthProvider
	UserRepository         domain.UserRepository
	RefreshTokenRepository domain.RefreshTokenRepository
}

var _ providersports.OAuthProviderService = OAuthGithubService{}

func NewOAuthGithubService(c domain.Config, p providersports.OAuthProvider, userRepository domain.UserRepository, refreshTokenRepository domain.RefreshTokenRepository) OAuthGithubService {
	return OAuthGithubService{
		Config:                 c,
		Provider:               p,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (s OAuthGithubService) GetAuthURL(redirectUri string) (string, error) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s",
		s.Config.Auth.Providers.GitHub.ClientID,
		"random_state_here", // TODO: generate proper state token and return it
	)
	return redirectURL, nil
}

func (s OAuthGithubService) ExchangeCode(code, state string) (http.Cookie, http.Cookie, error) {
	// todo: verify state, pass original state
	userInfos, err := s.Provider.GetUserInfos(code, state, "")
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	user, err := domain.GetOrCreateUserIfAllowed(s.UserRepository, userInfos, s.Config)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	if user.AuthMethod != "github" {
		return http.Cookie{}, http.Cookie{}, domain.ErrWrongAuthMethod
	}

	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := domain.GenerateTokensForUser(user, s.Config, s.RefreshTokenRepository)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessCookie := cookies.NewAccessCookie(accessToken, atExpiresAt, s.Config)
	refreshCookie := cookies.NewRefreshCookie(newRefreshToken, rtExpiresAt, s.Config)
	return accessCookie, refreshCookie, nil
}
