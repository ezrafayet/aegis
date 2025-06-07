package providers

import (
	"aegix/internal/domain"
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
