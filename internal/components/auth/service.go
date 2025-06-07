package auth

import (
	"aegix/internal/domain"
	"aegix/internal/repository"
	"aegix/pkg/cookies"
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
	if refreshToken != "" {
		_ = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	}
	return cookies.NewAccessCookieZero(s.Config), cookies.NewRefreshCookieZero(s.Config), nil
}

func (s AuthService) CheckAndRefreshToken(accessToken, refreshToken string) (http.Cookie, http.Cookie, error) {
	_, err := domain.ReadAccessTokenClaims(accessToken, s.Config)
	if err == nil {
		return http.Cookie{}, http.Cookie{}, domain.ErrValidAccessToken
	}
	if err.Error() != domain.ErrAccessTokenExpired.Error() {
		return http.Cookie{}, http.Cookie{}, err
	}
	refreshTokenObject, err := s.RefreshTokenRepository.GetRefreshTokenByToken(refreshToken)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}
	if refreshTokenObject.IsExpired() {
		return cookies.NewAccessCookieZero(s.Config), cookies.NewRefreshCookieZero(s.Config), domain.ErrRefreshTokenExpired
	}
	user, err := s.UserRepository.GetUserByID(refreshTokenObject.UserID)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}
	if user.IsDeleted() {
		return cookies.NewAccessCookieZero(s.Config), cookies.NewRefreshCookieZero(s.Config), domain.ErrUserDeleted
	}
	if user.IsBlocked() {
		return cookies.NewAccessCookieZero(s.Config), cookies.NewRefreshCookieZero(s.Config), domain.ErrUserBlocked
	}
	err = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		return cookies.NewAccessCookieZero(s.Config), cookies.NewRefreshCookieZero(s.Config), err
	}
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := domain.GenerateTokensForUser(user, s.Config, &s.RefreshTokenRepository)
	if err != nil {
		return cookies.NewAccessCookieZero(s.Config), cookies.NewRefreshCookieZero(s.Config), err
	}
	accessCookie := cookies.NewAccessCookie(accessToken, atExpiresAt, s.Config)
	refreshCookie := cookies.NewRefreshCookie(newRefreshToken, rtExpiresAt, s.Config)
	return accessCookie, refreshCookie, nil
}
