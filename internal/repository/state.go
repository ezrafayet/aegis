package repository

import (
	"othnx/internal/domain"

	"gorm.io/gorm"
)

type StateRepository struct {
	db *gorm.DB
}

var _ domain.StateRepository = &StateRepository{}

func NewStateRepository(db *gorm.DB) StateRepository {
	return StateRepository{db: db}
}

func (r *StateRepository) CreateState(state domain.State) error {
	return r.db.Create(&state).Error
}

func (r *StateRepository) GetAndDeleteState(value string) (domain.State, error) {
	var state domain.State
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("value = ?", value).First(&state).Error; err != nil {
			return err
		}
		if err := tx.Where("value = ?", value).Delete(&domain.State{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.State{}, err
	}
	return state, nil
}
