package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type RefreshToken struct {
	UserID    string    `json:"id" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"index;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null"`
	Token     string    `json:"token" gorm:"primaryKey;type:varchar(150);not null"`
	// DeviceFingerprint string   `json:"device_fingerprint" gorm:"type:varchar(150)"`
}

func (r RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}

func NewRefreshToken(user User, config Config) (RefreshToken, int64) {
	createdAt := time.Now()
	expiresAt := createdAt.AddDate(0, 0, config.JWT.RefreshTokenExpirationDays)
	return RefreshToken{
		UserID:    user.ID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
		Token:     uuid.New().String(),
	}, expiresAt.Unix()
}

type RefreshTokenRepository interface {
	CreateRefreshToken(refreshToken RefreshToken) error
	GetRefreshTokenByToken(token string) (RefreshToken, error)
	GetValidRefreshTokensByUserID(userID string) ([]RefreshToken, error)
	CleanExpiredTokens(userID string) error
	DeleteRefreshToken(token string) error
}


type CustomClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Metadata string `json:"metadata"`
}

func NewAccessToken(cClaims CustomClaims, config Config) (accessToken string, expiresAt int64, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	issuedAt := time.Now()
	secondsOfValidity := config.JWT.AccessTokenExpirationMin * 60
	expiresAt = issuedAt.Add(time.Second * time.Duration(secondsOfValidity)).Unix()
	claims := token.Claims.(jwt.MapClaims)
	claims["aud"] = config.App.Name
	claims["exp"] = expiresAt
	claims["issued_at"] = issuedAt.Unix()
	claims["iss"] = config.App.Name
	claims["user_id"] = cClaims.UserID
	// claims["roles"] = strings.Join(customClaims.Roles, ",") // todo
	claims["metadata"] = cClaims.Metadata
	tokenString, err := token.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		return "", -1, err
	}
	return tokenString, expiresAt, nil
}

func ReadAccessTokenClaims(accessToken string, config Config) (CustomClaims, error) {
	parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(config.JWT.Secret), nil
	})
	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors&jwt.ValidationErrorExpired != 0 {
				return CustomClaims{}, ErrAccessTokenExpired
			}
		}
		return CustomClaims{}, err
	}
	if !parsedToken.Valid {
		return CustomClaims{}, ErrInvalidToken
	}

	var customClaims CustomClaims

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		// /!\ This code can fail if the claims are not in the expected format
		customClaims.UserID = claims["user_id"].(string)
		// customClaims.Roles = strings.Split(claims["roles"].(string), " ")
		customClaims.Metadata = claims["metadata"].(string)
	}

	return customClaims, nil
}

func GenerateTokensForUser(user User, config Config, refreshTokenRepository RefreshTokenRepository) (accessToken string, atExpiresAt int64, refreshToken string, rtExpiresAt int64, err error) {
	validRefreshTokens, err := refreshTokenRepository.GetValidRefreshTokensByUserID(user.ID)
	if err != nil {
		return "", -1, "", -1, err
	}
	if len(validRefreshTokens) > 10 {
		return "", -1, "", -1, ErrTooManyRefreshTokens
	}

	_ = refreshTokenRepository.CleanExpiredTokens(user.ID)

	newRefreshToken, rtExpiresAt := NewRefreshToken(user, config)
	err = refreshTokenRepository.CreateRefreshToken(newRefreshToken)
	if err != nil {
		return "", -1, "", -1, err
	}

	accessToken, atExpiresAt, err = NewAccessToken(CustomClaims{
		UserID: user.ID,
		Metadata: user.Metadata,
	}, config)
	if err != nil {
		return "", -1, "", -1, err
	}

	return accessToken, atExpiresAt, newRefreshToken.Token, rtExpiresAt, nil
}
