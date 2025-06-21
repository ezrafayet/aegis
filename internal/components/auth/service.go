package auth

import (
	"net/http"
	"othnx/internal/domain"
	"othnx/internal/repository"
	"othnx/pkg/apperrors"
	"othnx/pkg/cookies"
)

type AuthServiceInterface interface {
	GetSession(accessToken string) (domain.Session, error)
	Logout(refreshToken string) (*http.Cookie, *http.Cookie, error)
	CheckAndRefreshToken(accessToken, refreshToken string) (*http.Cookie, *http.Cookie, error)
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

func (s AuthService) resetCookies(err error) (*http.Cookie, *http.Cookie, error) {
	ac, rc := cookies.NewAccessCookieZero(s.Config), cookies.NewRefreshCookieZero(s.Config)
	return &ac, &rc, err
}

func (s AuthService) Logout(refreshToken string) (*http.Cookie, *http.Cookie, error) {
	if refreshToken != "" {
		_ = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	}
	return s.resetCookies(nil)
}

func (s AuthService) CheckAndRefreshToken(accessToken, refreshToken string) (*http.Cookie, *http.Cookie, error) {
	_, err := domain.ReadAccessTokenClaims(accessToken, s.Config)
	if err == nil {
		return nil, nil, nil
	}
	if err.Error() != apperrors.ErrAccessTokenExpired.Error() {
		return s.resetCookies(err)
	}
	refreshTokenObject, err := s.RefreshTokenRepository.GetRefreshTokenByToken(refreshToken)
	if err != nil {
		return s.resetCookies(err)
	}
	if refreshTokenObject.IsExpired() {
		return s.resetCookies(apperrors.ErrRefreshTokenExpired)
	}
	// todo: check device id
	user, err := s.UserRepository.GetUserByID(refreshTokenObject.UserID)
	if err != nil {
		return s.resetCookies(err)
	}
	if user.IsDeleted() {
		return s.resetCookies(apperrors.ErrUserDeleted)
	}
	if user.IsBlocked() {
		return s.resetCookies(apperrors.ErrUserBlocked)
	}
	if s.Config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return s.resetCookies(apperrors.ErrEarlyAdoptersOnly)
	}

	err = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	if err != nil {
		return s.resetCookies(err)
	}
	// todo device-id: pass one, since one session per device is allowed
	accessToken, atExpiresAt, newRefreshToken, rtExpiresAt, err := domain.GenerateTokensForUser(user, "device-id", s.Config, &s.RefreshTokenRepository)
	if err != nil {
		return s.resetCookies(err)
	}
	accessCookie := cookies.NewAccessCookie(accessToken, atExpiresAt, s.Config)
	refreshCookie := cookies.NewRefreshCookie(newRefreshToken, rtExpiresAt, s.Config)
	return &accessCookie, &refreshCookie, nil
}
