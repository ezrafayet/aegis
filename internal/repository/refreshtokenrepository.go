package repository

import (
	"aegix/internal/domain"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

var _ domain.RefreshTokenRepository = &RefreshTokenRepository{}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(refreshToken domain.RefreshToken) error {
	return nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByToken(token string) (domain.RefreshToken, error) {
	return domain.RefreshToken{}, nil
}
