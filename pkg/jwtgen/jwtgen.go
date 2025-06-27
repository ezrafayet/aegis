package jwtgen

import (
	"othnx/internal/domain/entities"
	"othnx/pkg/apperrors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func Generate(cClaims entities.CustomClaims, config entities.Config, issuedAt time.Time) (accessToken string, expiresAt int64, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	secondsOfValidity := config.JWT.AccessTokenExpirationMin * 60
	expiresAt = issuedAt.Add(time.Second * time.Duration(secondsOfValidity)).Unix()
	claims := token.Claims.(jwt.MapClaims)
	claims["aud"] = config.App.Name
	claims["exp"] = expiresAt
	claims["issued_at"] = issuedAt.Unix()
	claims["iss"] = config.App.Name
	claims["user_id"] = cClaims.UserID
	claims["early_adopter"] = cClaims.EarlyAdopter
	claims["roles"] = strings.Join(cClaims.Roles, ",")
	claims["metadata"] = cClaims.Metadata
	tokenString, err := token.SignedString([]byte(config.JWT.Secret))
	if err != nil {
		return "", -1, err
	}
	return tokenString, expiresAt, nil
}

func ReadClaims(accessToken string, config entities.Config) (entities.CustomClaims, error) {
	parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrAccessTokenInvalid
		}
		return []byte(config.JWT.Secret), nil
	})
	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors&jwt.ValidationErrorExpired != 0 {
				return entities.CustomClaims{}, apperrors.ErrAccessTokenExpired
			}
		}
		return entities.CustomClaims{}, apperrors.ErrAccessTokenInvalid
	}
	if !parsedToken.Valid {
		return entities.CustomClaims{}, apperrors.ErrAccessTokenInvalid
	}

	var customClaims entities.CustomClaims

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		// /!\ This code can fail if the claims are not in the expected format
		customClaims.UserID = claims["user_id"].(string)
		customClaims.EarlyAdopter = claims["early_adopter"].(bool)
		customClaims.Roles = strings.Split(claims["roles"].(string), ",")
		customClaims.Metadata = claims["metadata"].(string)
	}

	return customClaims, nil
}
