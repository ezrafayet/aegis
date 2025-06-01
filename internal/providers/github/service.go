package github

import (
	"aegix/internal/domain"
	"aegix/internal/providers"
	"fmt"
)

type OAuthGithubService struct {
	Config   domain.Config
	Provider providers.OAuthProvider
	// DB
}

var _ providers.OAuthProviderService = OAuthGithubService{}

func NewOAuthGithubService(c domain.Config, p providers.OAuthProvider) OAuthGithubService {
	return OAuthGithubService{
		Config:   c,
		Provider: p,
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

func (s OAuthGithubService) ExchangeCode(code, state string) (string, error) {
	// todo: verify state, pass original state
	userInfos, err := s.Provider.GetUserInfos(code, state, "")
	if err != nil {
		return "", err
	}

	// save the user, generate token

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> userInfos", userInfos)

	return userInfos.Email, nil
}
