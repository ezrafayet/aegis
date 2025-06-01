package domain

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// Not persisted in a DB

type CustomClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
}

func NewAccessToken(customClaims CustomClaims, config Config) (accessToken string, expiresAt int64, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	issuedAt := time.Now()
	secondsOfValidity := config.JWT.AccessTokenExpirationMin * 60
	expiresAt = issuedAt.Add(time.Second * time.Duration(secondsOfValidity)).Unix()
	claims := token.Claims.(jwt.MapClaims)
	claims["aud"] = config.App.Name
	claims["exp"] = expiresAt
	claims["issued_at"] = issuedAt.Unix()
	claims["iss"] = config.App.Name
	claims["user_id"] = customClaims.UserID
	claims["roles"] = strings.Join(customClaims.Roles, ",")
	tokenString, err := token.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		return "", -1, err
	}
	return tokenString, expiresAt, nil
}
