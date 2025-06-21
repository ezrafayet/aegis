package repositories

import "othnx/internal/core/domain"

type StateRepository interface {
	CreateState(state domain.State) error
	GetAndDeleteState(value string) (domain.State, error)
}
