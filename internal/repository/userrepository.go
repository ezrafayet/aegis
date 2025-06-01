package repository

import (
	"aegix/internal/domain"
	"aegix/internal/providers"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var _ domain.UserRepository = &UserRepository{}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user domain.User) error {
	result := r.db.Model(&domain.User{}).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return domain.User{}, result.Error
	}
	if result.Error == gorm.ErrRecordNotFound {
		return domain.User{}, providers.ErrNoUser
	}
	return user, nil
}
