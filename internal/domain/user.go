package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         string     `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedAt  time.Time  `json:"created_at" gorm:"index;not null"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
	BlockedAt  *time.Time `json:"blocked_at" gorm:"index"`
	Name       string     `json:"name" gorm:"type:varchar(100);not null"`
	Avatar     string     `json:"avatar" gorm:"type:varchar(1000)"`
	Email      string     `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Metadata   string     `json:"metadata" gorm:"type:varchar(1000)"`
	Roles      []string   `json:"roles" gorm:"type:jsonb"`
	AuthMethod string     `json:"auth_method" gorm:"type:varchar(20);not null"`
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
