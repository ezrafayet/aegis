package jwtgen

import (
	"errors"
	"time"

	"auth302/pkg/apperrors"

	"github.com/golang-jwt/jwt/v5"
)

func Generate(
	cClaims map[string]any,
	issuedAt time.Time,
	accessTokenExpirationMin int,
	appName string,
	secret string,
) (accessToken string, expiresAt int64, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	secondsOfValidity := accessTokenExpirationMin * 60
	expiresAt = issuedAt.Add(time.Second * time.Duration(secondsOfValidity)).Unix()
	claims := token.Claims.(jwt.MapClaims)
	claims["aud"] = appName
	claims["exp"] = expiresAt
	claims["issued_at"] = issuedAt.Unix()
	claims["iss"] = appName
	for key, value := range cClaims {
		claims[key] = value
	}
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", -1, err
	}
	return tokenString, expiresAt, nil
}

func ReadClaims(accessToken string, secret string) (map[string]any, error) {
	parsedToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrAccessTokenInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return map[string]any{}, apperrors.ErrAccessTokenExpired
		}
		return map[string]any{}, apperrors.ErrAccessTokenInvalid
	}
	if !parsedToken.Valid {
		return map[string]any{}, apperrors.ErrAccessTokenInvalid
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return map[string]any{}, apperrors.ErrAccessTokenInvalid
	}
	return claims, nil
}
