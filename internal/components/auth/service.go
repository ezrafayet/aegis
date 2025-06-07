package auth

import (
	"aegix/internal/domain"
	"aegix/internal/repository"
	"aegix/pkg/cookies"
	"errors"
	"net/http"
)

type AuthServiceInterface interface {
	GetSession(accessToken string) (domain.Session, error)
	Logout(refreshToken string) (http.Cookie, http.Cookie, error)
	CheckAndRefreshToken(accessToken, refreshToken string) (http.Cookie, http.Cookie, error)
}

type AuthService struct {
	Config                 domain.Config
	RefreshTokenRepository repository.RefreshTokenRepository
	UserRepository         repository.UserRepository
}

var _ AuthServiceInterface = &AuthService{}

func NewAuthService(c domain.Config, r repository.RefreshTokenRepository, u repository.UserRepository) AuthService {
	return AuthService{
		Config:                 c,
		RefreshTokenRepository: r,
		UserRepository:         u,
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
	refreshCookie := cookies.NewRefreshCookie("", 0, s.Config)
	accessCookie := cookies.NewAccessCookie("", 0, s.Config)
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
		refreshCookie := cookies.NewRefreshCookie("", 0, s.Config)
		accessCookie := cookies.NewAccessCookie("", 0, s.Config)
		return refreshCookie, accessCookie, errors.New("refresh_token_expired")
	}

	user, err := s.UserRepository.GetUserByID(refreshTokenObject.UserID)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	if user.IsBlocked() {
		refreshCookie := cookies.NewRefreshCookie("", 0, s.Config)
		accessCookie := cookies.NewAccessCookie("", 0, s.Config)
		return refreshCookie, accessCookie, errors.New("user_blocked")
	}

	if user.IsDeleted() {
		refreshCookie := cookies.NewRefreshCookie("", 0, s.Config)
		accessCookie := cookies.NewAccessCookie("", 0, s.Config)
		return refreshCookie, accessCookie, errors.New("user_deleted")
	}

	err = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := domain.GenerateTokensForUser(user, s.Config, &s.RefreshTokenRepository)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessCookie := cookies.NewAccessCookie(accessToken, atExpiresAt, s.Config)
	refreshCookie := cookies.NewRefreshCookie(newRefreshToken, rtExpiresAt, s.Config)
	return accessCookie, refreshCookie, nil
}
