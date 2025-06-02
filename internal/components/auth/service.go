package auth

import (
	"aegix/internal/components/cookies"
	"aegix/internal/domain"
	"aegix/internal/providers"
	"aegix/internal/repository"
	"errors"
	"net/http"
)

type AuthServiceInterface interface {
	GetSession(accessToken string) (domain.Session, error)
	Logout(refreshToken string) (http.Cookie, http.Cookie, error)
	CheckAndRefreshToken(accessToken, refreshToken string) (http.Cookie, http.Cookie, error)
}

type AuthService struct {
	Config domain.Config
	RefreshTokenRepository repository.RefreshTokenRepository
	UserRepository repository.UserRepository
	CookieBuilder          cookies.CookieBuilderMethods
}

var _ AuthServiceInterface = &AuthService{}

func NewAuthService(c domain.Config, r repository.RefreshTokenRepository, u repository.UserRepository) AuthService {
	return AuthService{
		Config: c,
		RefreshTokenRepository: r,
		UserRepository: u,
		CookieBuilder: cookies.NewCookieBuilder(c), // how bad is that?
	}
}

func (s AuthService) GetSession(accessToken string) (domain.Session, error) {
	customClaims, err := domain.ReadAccessTokenClaims(accessToken, s.Config)
	if err != nil {
		return domain.Session{}, err
	}

	return domain.Session{
		CustomClaims: customClaims,
	}, nil
}

func (s AuthService) Logout(refreshToken string) (http.Cookie, http.Cookie, error) {
	refreshCookie := s.CookieBuilder.NewRefreshCookie("", 0, true)
	accessCookie := s.CookieBuilder.NewAccessCookie("", 0, true)
	err := s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		return refreshCookie, accessCookie, err
	}
	return refreshCookie, accessCookie, nil
}

func (s AuthService) CheckAndRefreshToken(accessToken, refreshToken string) (http.Cookie, http.Cookie, error) {
	_, err := domain.ReadAccessTokenClaims(accessToken, s.Config)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	refreshTokenObject, err := s.RefreshTokenRepository.GetRefreshTokenByToken(refreshToken)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	if refreshTokenObject.IsExpired() {
		refreshCookie := s.CookieBuilder.NewRefreshCookie("", 0, true)
		accessCookie := s.CookieBuilder.NewAccessCookie("", 0, true)
		return refreshCookie, accessCookie, errors.New("refresh_token_expired")
	}

	user, err := s.UserRepository.GetUserByID(refreshTokenObject.UserID)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	if user.IsBlocked() {
		refreshCookie := s.CookieBuilder.NewRefreshCookie("", 0, true)
		accessCookie := s.CookieBuilder.NewAccessCookie("", 0, true)
		return refreshCookie, accessCookie, errors.New("user_blocked")
	}

	if user.IsDeleted() {
		refreshCookie := s.CookieBuilder.NewRefreshCookie("", 0, true)
		accessCookie := s.CookieBuilder.NewAccessCookie("", 0, true)
		return refreshCookie, accessCookie, errors.New("user_deleted")
	}

	err = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	// duplicated code btw
	// arbitrary naive check, will replace with device fingerprints
	validRefreshTokens, err := s.RefreshTokenRepository.GetValidRefreshTokensByUserID(user.ID)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}
	if len(validRefreshTokens) > 10 {
		return http.Cookie{}, http.Cookie{}, providers.ErrTooManyRefreshTokens
	}

	_ = s.RefreshTokenRepository.CleanExpiredTokens(user.ID)

	newRefreshToken, rtExpiresAt := domain.NewRefreshToken(user, s.Config)
	err = s.RefreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessToken, atExpiresAt, err := domain.NewAccessToken(user, s.Config)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessCookie := s.CookieBuilder.NewAccessCookie(accessToken, atExpiresAt, true)
	refreshCookie := s.CookieBuilder.NewRefreshCookie(newRefreshToken.Token, rtExpiresAt, true)
	return accessCookie, refreshCookie, nil
}
