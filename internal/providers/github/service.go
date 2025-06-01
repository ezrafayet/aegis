package github

import (
	"aegix/internal/domain"
	"aegix/internal/providers"
	"fmt"
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

func (s OAuthGithubService) ExchangeCode(code, state string) (string, error) {
	// todo: verify state, pass original state
	userInfos, err := s.Provider.GetUserInfos(code, state, "")
	if err != nil {
		return "", err
	}

	user, err := s.UserRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != providers.ErrNoUser.Error() {
		return "", err
	}

	if err != nil && err.Error() == providers.ErrNoUser.Error() {
		user = domain.NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, "github")
		err = s.UserRepository.CreateUser(user)
		if err != nil {
			return "", err
		}
	}

	if user.IsBlocked() {
		return "", providers.ErrUserBlocked
	}

	if user.IsDeleted() {
		return "", providers.ErrUserDeleted
	}

	if user.AuthMethod != "github" {
		return "", providers.ErrWrongAuthMethod
	}

	var refreshToken domain.RefreshToken

	validRefreshTokens, err := s.RefreshTokenRepository.GetValidRefreshTokensByUserID(user.ID)
	if err != nil {
		return "", err
	}

	if len(validRefreshTokens) == 0 {
		refreshToken = domain.NewRefreshToken(user.ID, s.Config.JWT.RefreshTokenExpirationDays)
		err = s.RefreshTokenRepository.CreateRefreshToken(refreshToken)
		if err != nil {
			return "", err
		}
	} else {
		refreshToken = validRefreshTokens[0]
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>> refreshToken", refreshToken)
	// create jwt and cookie

	return userInfos.Email, nil
}
