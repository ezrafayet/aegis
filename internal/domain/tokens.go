package domain

func GenerateTokensForUser(user User, config Config, refreshTokenRepository RefreshTokenRepository) (accessToken string, atExpiresAt int64, refreshToken string, rtExpiresAt int64, err error) {
	validRefreshTokens, err := refreshTokenRepository.GetValidRefreshTokensByUserID(user.ID)
	if err != nil {
		return "", -1, "", -1, err
	}
	if len(validRefreshTokens) > 10 {
		return "", -1, "", -1, ErrTooManyRefreshTokens
	}

	_ = refreshTokenRepository.CleanExpiredTokens(user.ID)

	newRefreshToken, rtExpiresAt := NewRefreshToken(user, config)
	err = refreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return "", -1, "", -1, err
	}

	accessToken, atExpiresAt, err = NewAccessToken(CustomClaims{
		UserID: user.ID,
		Metadata: user.Metadata,
	}, config)
	if err != nil {
		return "", -1, "", -1, err
	}

	return accessToken, atExpiresAt, newRefreshToken.Token, rtExpiresAt, nil
}
