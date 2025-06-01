package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	BlockedAt *time.Time `json:"blocked_at"`
	Username  string     `json:"username"`
	Avatar    string     `json:"avatar"`
	Email     string     `json:"email"`
	Metadata  string     `json:"metadata"`
	Roles     []string   `json:"roles"`
	AuthMethod string    `json:"auth_method"`
}

func NewUser(username, avatar, email string, roles []string, authMethod string) User {
	return User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		DeletedAt: nil,
		BlockedAt: nil,
		Username:  username,
		Email:     email,
		Metadata:  "{}",
		Roles:     roles,
		AuthMethod: authMethod,
	}
}

type UserRepository interface {
	CreateUser(user User) error
	// GetUserByID(userID string) (User, error)
	GetUserByEmail(email string) (User, error)
	// SoftDeleteUser(userID string) error
	// HardDeleteUser(userID string) error
	// BlockUser(userID string) error
	// UnblockUser(userID string) error
}
