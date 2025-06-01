package github

import (
	"aegix/internal/domain"
	"aegix/internal/providers"
	"fmt"
	"net/http"
	"time"
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

	user, err := s.UserRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != providers.ErrNoUser.Error() {
		return http.Cookie{}, http.Cookie{}, err
	}

	if err != nil && err.Error() == providers.ErrNoUser.Error() {
		user = domain.NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, "github")
		err = s.UserRepository.CreateUser(user)
		if err != nil {
			return http.Cookie{}, http.Cookie{}, err
		}
	}

	if user.IsBlocked() {
		return http.Cookie{}, http.Cookie{}, providers.ErrUserBlocked
	}

	if user.IsDeleted() {
		return http.Cookie{}, http.Cookie{}, providers.ErrUserDeleted
	}

	if user.AuthMethod != "github" {
		return http.Cookie{}, http.Cookie{}, providers.ErrWrongAuthMethod
	}

	validRefreshTokens, err := s.RefreshTokenRepository.GetValidRefreshTokensByUserID(user.ID)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	// arbitrary naive check, will replace with device fingerprints
	if len(validRefreshTokens) > 10 {
		return http.Cookie{}, http.Cookie{}, providers.ErrTooManyRefreshTokens
	}

	_ = s.RefreshTokenRepository.CleanExpiredTokens(user.ID)

	refreshToken, rtExpiresAt := domain.NewRefreshToken(user.ID, s.Config)
	err = s.RefreshTokenRepository.CreateRefreshToken(refreshToken)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessToken, atExpiresAt, err := domain.NewAccessToken(domain.CustomClaims{
		UserID: user.ID,
		Roles:  []string{}, // replace with user roles
	}, s.Config)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Domain:   s.Config.App.URL,
		Value:    accessToken,
		Expires:  time.Unix(atExpiresAt, 0),
		HttpOnly: true,
		Secure:   false, // change
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Domain:   s.Config.App.URL,
		Value:    refreshToken.Token,
		Expires:  time.Unix(rtExpiresAt, 0),
		HttpOnly: true,
		Secure:   false, // change
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	return accessCookie, refreshCookie, nil
}
