package usecases

import (
	"othnx/internal/domain/entities"
	"othnx/internal/domain/ports/secondary_ports"
	"othnx/pkg/apperrors"
	"othnx/pkg/fingerprint"
	"othnx/pkg/jwtgen"
	"time"
)

func GetOrCreateUserIfAllowed(userRepository secondaryports.UserRepository, userInfos *entities.UserInfos, config entities.Config) (entities.User, error) {
	nameExists, err := userRepository.DoesNameExist(userInfos.Name)
	if err != nil {
		return entities.User{}, err
	}
	if nameExists {
		return entities.User{}, apperrors.ErrNameAlreadyExists
	}
	user, err := userRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != apperrors.ErrNoUser.Error() {
		return entities.User{}, err
	}
	if err != nil && err.Error() == apperrors.ErrNoUser.Error() {
		user, err = entities.NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, "github")
		if err != nil {
			return entities.User{}, err
		}
		err = userRepository.CreateUser(user, []entities.Role{entities.NewRole(user.ID, entities.RoleUser)})
		if err != nil {
			return entities.User{}, err
		}
	}
	if user.IsDeleted() {
		return entities.User{}, apperrors.ErrUserDeleted
	}
	if user.IsBlocked() {
		return entities.User{}, apperrors.ErrUserBlocked
	}
	if config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return entities.User{}, apperrors.ErrEarlyAdoptersOnly
	}
	return user, nil
}

func GenerateTokensForUser(user entities.User, deviceID string, config entities.Config, refreshTokenRepository secondaryports.RefreshTokenRepository) (accessToken string, atExpiresAt int64, refreshToken string, rtExpiresAt int64, err error) {
	deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint(deviceID)
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

	newRefreshToken, rtExpiresAt, err := entities.NewRefreshToken(user, deviceFingerprint, config)
	if err != nil {
		return "", -1, "", -1, err
	}
	err = refreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return "", -1, "", -1, err
	}

	accessToken, atExpiresAt, err = jwtgen.Generate(entities.CustomClaims{
		UserID:       user.ID,
		EarlyAdopter: user.EarlyAdopter,
		Metadata:     user.Metadata,
		Roles:        user.RolesValues(),
	}, config, time.Now())
	if err != nil {
		return "", -1, "", -1, err
	}

	return accessToken, atExpiresAt, newRefreshToken.Token, rtExpiresAt, nil
}
