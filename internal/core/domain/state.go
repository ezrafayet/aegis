package domain

import "time"

type State struct {
	Value     string    `json:"value" gorm:"type:char(32);index;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null"`
}

func (s State) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func NewState(value string) State {
	return State{
		Value:     value,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
}
