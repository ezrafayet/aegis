package github

import (
	"aegix/internal/components/providers"
	"aegix/internal/domain"
	"aegix/pkg/cookies"
	"fmt"
	"net/http"
)

type OAuthGithubService struct {
	Config                 domain.Config
	Provider               providers.OAuthProvider
	UserRepository         domain.UserRepository
	RefreshTokenRepository domain.RefreshTokenRepository
}

var _ providers.OAuthProviderService = OAuthGithubService{}

func NewOAuthGithubService(c domain.Config, p providers.OAuthProvider, userRepository domain.UserRepository, refreshTokenRepository domain.RefreshTokenRepository) OAuthGithubService {
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

	user, err := providers.GetOrCreateUserIfAllowed(s.UserRepository, userInfos)
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

	accessCookie := cookies.NewAccessCookie(accessToken, atExpiresAt, true, s.Config)
	refreshCookie := cookies.NewRefreshCookie(newRefreshToken, rtExpiresAt, true, s.Config)
	return accessCookie, refreshCookie, nil
}
