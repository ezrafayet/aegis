package providers

import (
	"aegix/internal/domain"
	"aegix/pkg/cookies"
	"net/http"
)

func GetOrCreateUserIfAllowed(userRepository domain.UserRepository, userInfos *OAuthUser) (domain.User, error) {
	user, err := userRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != domain.ErrNoUser.Error() {
		return domain.User{}, err
	}

	if err != nil && err.Error() == domain.ErrNoUser.Error() {
		user = domain.NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, "github")
		err = userRepository.CreateUser(user)
		if err != nil {
			return domain.User{}, err
		}
	}

	if user.IsDeleted() {
		return domain.User{}, domain.ErrUserDeleted
	}

	if user.IsBlocked() {
		return domain.User{}, domain.ErrUserBlocked
	}

	return user, nil
}

func GetTokensForUser(user domain.User, refreshTokenRepository domain.RefreshTokenRepository, config domain.Config) (http.Cookie, http.Cookie, error) {
	// arbitrary naive check, will replace with device fingerprints
	validRefreshTokens, err := refreshTokenRepository.GetValidRefreshTokensByUserID(user.ID)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}
	if len(validRefreshTokens) > 10 {
		return http.Cookie{}, http.Cookie{}, domain.ErrTooManyRefreshTokens
	}

	_ = refreshTokenRepository.CleanExpiredTokens(user.ID)

	refreshToken, rtExpiresAt := domain.NewRefreshToken(user, config)
	err = refreshTokenRepository.CreateRefreshToken(refreshToken)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessToken, atExpiresAt, err := domain.NewAccessToken(user, config)
	if err != nil {
		return http.Cookie{}, http.Cookie{}, err
	}

	accessCookie := cookies.NewAccessCookie(accessToken, atExpiresAt, true, config)
	refreshCookie := cookies.NewRefreshCookie(refreshToken.Token, rtExpiresAt, true, config)

	return accessCookie, refreshCookie, nil
}
