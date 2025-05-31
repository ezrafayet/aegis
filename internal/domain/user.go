package domain

import "time"

type UserId string

type User struct {
	ID        UserId     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	BlockedAt *time.Time `json:"blocked_at"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Metadata  string     `json:"metadata"`
	Roles     []string   `json:"roles"`
}
