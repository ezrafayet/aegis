package use_cases

func GetOrCreateUserIfAllowed(userRepository repositories.UserRepository, userInfos *UserInfos, config Config) (User, error) {
	nameExists, err := userRepository.DoesNameExist(userInfos.Name)
	if err != nil {
		return User{}, err
	}
	if nameExists {
		return User{}, apperrors.ErrNameAlreadyExists
	}
	user, err := userRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != apperrors.ErrNoUser.Error() {
		return User{}, err
	}
	if err != nil && err.Error() == apperrors.ErrNoUser.Error() {
		user, err = NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, "github")
		if err != nil {
			return User{}, err
		}
		err = userRepository.CreateUser(user, []Role{NewRole(user.ID, "user")})
		if err != nil {
			return User{}, err
		}
	}
	if user.IsDeleted() {
		return User{}, apperrors.ErrUserDeleted
	}
	if user.IsBlocked() {
		return User{}, apperrors.ErrUserBlocked
	}
	if config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return User{}, apperrors.ErrEarlyAdoptersOnly
	}
	return user, nil
}

func GenerateTokensForUser(user User, deviceID string, config Config, refreshTokenRepository RefreshTokenRepository) (accessToken string, atExpiresAt int64, refreshToken string, rtExpiresAt int64, err error) {
	deviceFingerprint, err := GenerateDeviceFingerprint(deviceID)
	if err != nil {
		return "", -1, "", -1, err
	}

	err = refreshTokenRepository.DeleteRefreshTokenByDeviceFingerprint(user.ID, deviceFingerprint)
	if err != nil {
		return "", -1, "", -1, err
	}

	validRefreshTokens, err := refreshTokenRepository.CountValidRefreshTokensForUser(user.ID)
	if err != nil {
		return "", -1, "", -1, err
	}

	if validRefreshTokens >= 5 {
		return "", -1, "", -1, apperrors.ErrTooManyRefreshTokens
	}

	_ = refreshTokenRepository.CleanExpiredTokens(user.ID)

	newRefreshToken, rtExpiresAt, err := NewRefreshToken(user, deviceFingerprint, config)
	if err != nil {
		return "", -1, "", -1, err
	}
	err = refreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return "", -1, "", -1, err
	}

	accessToken, atExpiresAt, err = NewAccessToken(CustomClaims{
		UserID:       user.ID,
		EarlyAdopter: user.EarlyAdopter,
		Metadata:     user.Metadata,
		RolesValues:  user.RolesValues(),
	}, config, time.Now())
	if err != nil {
		return "", -1, "", -1, err
	}

	return accessToken, atExpiresAt, newRefreshToken.Token, rtExpiresAt, nil
}
