package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string     `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedAt    time.Time  `json:"created_at" gorm:"index;not null"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
	BlockedAt    *time.Time `json:"blocked_at" gorm:"index"`
	EarlyAdopter bool       `json:"early_adopter" gorm:"index"`
	Name         string     `json:"name" gorm:"type:varchar(100);not null"`
	// NameFingerprint string `json:"name_fingerprint" gorm:"type:varchar(100);uniqueIndex;not null"`
	AvatarURL string `json:"avatar_url" gorm:"type:varchar(1000)"`
	Email     string `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Metadata  string `json:"metadata" gorm:"type:varchar(1000)"`
	// Roles      postgres.StringArray   `json:"roles" gorm:"type:text[]"`
	AuthMethod string `json:"auth_method" gorm:"type:varchar(20);not null"`
}

func (u User) IsEarlyAdopter() bool {
	return u.EarlyAdopter
}

func (u User) IsBlocked() bool {
	return u.BlockedAt != nil
}

func (u User) IsDeleted() bool {
	return u.DeletedAt != nil
}

func NewUser(name, avatar, email string, authMethod string) User {
	return User{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now(),
		DeletedAt:  nil,
		BlockedAt:  nil,
		EarlyAdopter: false,
		Name:       name,
		// NameFingerprint: compute nameFingerprint,
		AvatarURL: avatar,
		Email:     email,
		Metadata:  "{}",
		// Roles:      roles,
		AuthMethod: authMethod,
	}
}

type UserRepository interface {
	CreateUser(user User) error
	GetUserByID(userID string) (User, error)
	GetUserByEmail(email string) (User, error)
	// SoftDeleteUser(userID string) error
	// HardDeleteUser(userID string) error
	// BlockUser(userID string) error
	// UnblockUser(userID string) error
}

type UserInfos struct {
	Name   string
	Email  string
	Avatar string
}

func GetOrCreateUserIfAllowed(userRepository UserRepository, userInfos *UserInfos, config Config) (User, error) {
	user, err := userRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != ErrNoUser.Error() {
		return User{}, err
	}

	if err != nil && err.Error() == ErrNoUser.Error() {
		user = NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, "github")
		err = userRepository.CreateUser(user)
		if err != nil {
			return User{}, err
		}
	}

	if user.IsDeleted() {
		return User{}, ErrUserDeleted
	}

	if user.IsBlocked() {
		return User{}, ErrUserBlocked
	}

	if config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return User{}, ErrEarlyAdoptersOnly
	}

	return user, nil
}
