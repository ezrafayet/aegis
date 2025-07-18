package entities

import (
	"aegis/pkg/fingerprint"
	"aegis/pkg/uidgen"
	"time"
)

type User struct {
	ID              string     `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedAt       time.Time  `json:"created_at" gorm:"index;not null"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
	BlockedAt       *time.Time `json:"blocked_at" gorm:"index"`
	EarlyAdopter    bool       `json:"early_adopter" gorm:"index;default:false"`
	Name            string     `json:"name" gorm:"type:varchar(100);not null"`
	NameFingerprint string     `json:"name_fingerprint" gorm:"type:char(32);index;not null"` // not unique
	AvatarURL       string     `json:"avatar_url" gorm:"type:varchar(1024)"`
	Email           string     `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	MetadataPublic  string     `json:"metadata_public" gorm:"type:varchar(1024);not null"`
	AuthMethod      string     `json:"auth_method" gorm:"type:varchar(16);not null"`
	// relations
	Roles         []Role         `json:"roles" gorm:"foreignKey:UserID;references:ID"`
	RefreshTokens []RefreshToken `json:"refresh_tokens" gorm:"foreignKey:UserID;references:ID"`
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

func NewUser(name, avatar, email string, authMethod string) (User, error) {
	nameFingerprint, err := fingerprint.GenerateNameFingerprint(name)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:              uidgen.Generate(),
		CreatedAt:       time.Now(),
		DeletedAt:       nil,
		BlockedAt:       nil,
		EarlyAdopter:    false,
		Name:            name,
		NameFingerprint: nameFingerprint,
		AvatarURL:       avatar,
		Email:           email,
		MetadataPublic:  "{}",
		AuthMethod:      authMethod,
	}, nil
}
