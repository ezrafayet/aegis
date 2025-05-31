package domain

import "time"

type RefreshToken struct {
	ID        UserId    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
}
