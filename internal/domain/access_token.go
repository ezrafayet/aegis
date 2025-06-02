package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// Not persisted in a DB

type CustomClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	Metadata string `json:"metadata"`
}

func NewAccessToken(user User, config Config) (accessToken string, expiresAt int64, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	issuedAt := time.Now()
	secondsOfValidity := config.JWT.AccessTokenExpirationMin * 60
	expiresAt = issuedAt.Add(time.Second * time.Duration(secondsOfValidity)).Unix()
	claims := token.Claims.(jwt.MapClaims)
	claims["aud"] = config.App.Name
	claims["exp"] = expiresAt
	claims["issued_at"] = issuedAt.Unix()
	claims["iss"] = config.App.Name
	claims["user_id"] = user.ID
	// claims["roles"] = strings.Join(customClaims.Roles, ",") // todo
	claims["metadata"] = user.Metadata
	tokenString, err := token.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		return "", -1, err
	}
	return tokenString, expiresAt, nil
}

func ReadAccessTokenClaims(accessToken string, config Config) (CustomClaims, error) {
	parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid_token")
		}
		return []byte(config.JWT.Secret), nil
	})
	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors&jwt.ValidationErrorExpired != 0 {
				return CustomClaims{}, errors.New("access_token_expired")
			}
		}
		return CustomClaims{}, err
	}
	if !parsedToken.Valid {
		return CustomClaims{}, errors.New("invalid_token")
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
