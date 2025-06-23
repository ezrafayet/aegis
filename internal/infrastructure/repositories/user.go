package repositories

import (
	"othnx/internal/domain"
	"othnx/pkg/apperrors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var _ domain.UserRepository = &UserRepository{}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user domain.User, roles []domain.Role) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.User{}).Create(&user).Error; err != nil {
			return err
		}
		for _, role := range roles {
			if err := tx.Model(&domain.Role{}).Create(&role).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).Preload("Roles").First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return domain.User{}, result.Error
	}
	if result.Error == gorm.ErrRecordNotFound {
		return domain.User{}, apperrors.ErrNoUser
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(userID string) (domain.User, error) {
	var user domain.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepository) DoesNameExist(nameFingerprint string) (bool, error) {
	var user domain.User
	result := r.db.Where("name_fingerprint = ?", nameFingerprint).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}
	return result.Error == nil, nil
}
