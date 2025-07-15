package repositories

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"aegis/pkg/apperrors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var _ secondary.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user entities.User, roles []entities.Role) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.User{}).Create(&user).Error; err != nil {
			return err
		}
		for _, role := range roles {
			if err := tx.Model(&entities.Role{}).Create(&role).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (entities.User, error) {
	var user entities.User
	result := r.db.Where("email = ?", email).Preload("Roles").First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return entities.User{}, result.Error
	}
	if result.Error == gorm.ErrRecordNotFound {
		return entities.User{}, apperrors.ErrNoUser
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(userID string) (entities.User, error) {
	var user entities.User
	result := r.db.Where("id = ?", userID).Preload("Roles").First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return entities.User{}, result.Error
	}
	if result.Error == gorm.ErrRecordNotFound {
		return entities.User{}, apperrors.ErrNoUser
	}
	return user, nil
}

func (r *UserRepository) DoesNameExist(nameFingerprint string) (bool, error) {
	var user entities.User
	result := r.db.Where("name_fingerprint = ?", nameFingerprint).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.Error == nil, nil
}
