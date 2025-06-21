package domain

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"othnx/internal/core/repositories"
	"othnx/pkg/apperrors"
	"othnx/pkg/uidgen"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type User struct {
	ID              string     `json:"id" gorm:"primaryKey;type:uuid"`
	CreatedAt       time.Time  `json:"created_at" gorm:"index;not null"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
	BlockedAt       *time.Time `json:"blocked_at" gorm:"index"`
	EarlyAdopter    bool       `json:"early_adopter" gorm:"index;default:false"`
	Name            string     `json:"name" gorm:"type:varchar(100);not null"`
	NameFingerprint string     `json:"name_fingerprint" gorm:"type:char(32);uniqueIndex;not null"`
	AvatarURL       string     `json:"avatar_url" gorm:"type:varchar(1024)"`
	Email           string     `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Metadata        string     `json:"metadata" gorm:"type:varchar(1024);not null"`
	AuthMethod      string     `json:"auth_method" gorm:"type:varchar(16);not null"`

	// relations
	Roles []Role `json:"roles" gorm:"foreignKey:UserID;references:ID"`
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

func (u User) RolesValues() []string {
	values := make([]string, len(u.Roles))
	for i, role := range u.Roles {
		values[i] = role.Value
	}
	return values
}

func NewUser(name, avatar, email string, authMethod string) (User, error) {
	nameFingerprint, err := GenerateNameFingerprint(name)
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
		Metadata:        "{}",
		AuthMethod:      authMethod,
	}, nil
}

func GenerateNameFingerprint(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", errors.New("empty_name")
	}
	normalized := strings.ToLower(trimmed)
	transformer := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, err := transform.String(transformer, normalized)
	if err != nil {
		result = normalized
	}
	result = strings.Join(strings.Fields(result), " ")
	hash := md5.Sum([]byte(result))
	return hex.EncodeToString(hash[:]), nil
}

// UserInfos is what is returned by the providers (GitHub, Google, etc.)
type UserInfos struct {
	Name   string
	Email  string
	Avatar string
}

// todo: move this business logic somewhere
func GetOrCreateUserIfAllowed(userRepository repositories.UserRepository, userInfos *UserInfos, config Config) (User, error) {
	nameExists, err := userRepository.DoesNameExist(userInfos.Name)
	if err != nil {
		return User{}, err
	}
	if nameExists {
		return User{}, apperrors.ErrNameAlreadyExists
	}
	user, err := userRepository.GetUserByEmail(userInfos.Email)
	if err != nil && err.Error() != apperrors.ErrNoUser.Error() {
		return User{}, err
	}
	if err != nil && err.Error() == apperrors.ErrNoUser.Error() {
		user, err = NewUser(userInfos.Name, userInfos.Avatar, userInfos.Email, "github")
		if err != nil {
			return User{}, err
		}
		err = userRepository.CreateUser(user, []Role{NewRole(user.ID, "user")})
		if err != nil {
			return User{}, err
		}
	}
	if user.IsDeleted() {
		return User{}, apperrors.ErrUserDeleted
	}
	if user.IsBlocked() {
		return User{}, apperrors.ErrUserBlocked
	}
	if config.App.EarlyAdoptersOnly && !user.IsEarlyAdopter() {
		return User{}, apperrors.ErrEarlyAdoptersOnly
	}
	return user, nil
}
