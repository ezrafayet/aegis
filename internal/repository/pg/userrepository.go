package pgrepository

import (
	"aegix/internal/domain"

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
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	return domain.User{}, nil
}
