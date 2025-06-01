package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         string     `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	BlockedAt  *time.Time `json:"blocked_at"`
	Name       string     `json:"name"`
	Avatar     string     `json:"avatar"`
	Email      string     `json:"email"`
	Metadata   string     `json:"metadata"`
	Roles      []string   `json:"roles"`
	AuthMethod string     `json:"auth_method"`
}

func (u User) IsBlocked() bool {
	return u.BlockedAt != nil
}

func (u User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func NewUser(name, avatar, email string, roles []string, authMethod string) User {
	return User{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now(),
		DeletedAt:  nil,
		BlockedAt:  nil,
		Name:       name,
		Avatar:     avatar,
		Email:      email,
		Metadata:   "{}",
		Roles:      roles,
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
