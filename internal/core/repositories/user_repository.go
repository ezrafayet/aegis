package repositories

import "othnx/internal/core/domain"

type UserRepository interface {
	CreateUser(user domain.User, roles []domain.Role) error
	GetUserByID(userID string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	DoesNameExist(nameFingerprint string) (bool, error)
}
