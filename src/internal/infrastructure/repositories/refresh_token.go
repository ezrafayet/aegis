package repositories

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"
	"aegis/pkg/apperrors"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

var _ secondary.RefreshTokenRepository = (*RefreshTokenRepository)(nil)

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(refreshToken entities.RefreshToken) error {
	result := r.db.Model(&entities.RefreshToken{}).Create(&refreshToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByToken(token string) (entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	result := r.db.Model(&entities.RefreshToken{}).Where("token = ?", token).First(&refreshToken)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return entities.RefreshToken{}, result.Error
	}
	if result.Error == gorm.ErrRecordNotFound {
		return entities.RefreshToken{}, apperrors.ErrRefreshTokenInvalid
	}
	return refreshToken, nil
}

func (r *RefreshTokenRepository) CountValidRefreshTokensForUser(userID string) (int, error) {
	var count int64
	result := r.db.Model(&entities.RefreshToken{}).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}

func (r *RefreshTokenRepository) CleanExpiredTokens(userID string) error {
	result := r.db.Model(&entities.RefreshToken{}).Where("user_id = ? AND expires_at < ?", userID, time.Now()).Delete(&entities.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	result := r.db.Model(&entities.RefreshToken{}).Where("token = ?", token).Delete(&entities.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteRefreshTokenByDeviceFingerprint(userID, deviceFingerprint string) error {
	result := r.db.Model(&entities.RefreshToken{}).Where("user_id = ? AND device_fingerprint = ?", userID, deviceFingerprint).Delete(&entities.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
