package jwtgen

import (
	"othnx/internal/domain/entities"
	"othnx/internal/infrastructure/config"
	"othnx/pkg/apperrors"
	"testing"
	"time"
)

func TestAccessToken(t *testing.T) {
	t.Run("should create a new access token", func(t *testing.T) {
		token, expires_at, err := Generate(entities.CustomClaims{
			UserID:   "123",
			Roles:    []string{""},
			Metadata: "{foo:bar}",
		}, config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJlYXJseV9hZG9wdGVyIjpmYWxzZSwiZXhwIjoxNjcyNTMyMTAwLCJpc3MiOiIiLCJpc3N1ZWRfYXQiOjE2NzI1MzEyMDAsIm1ldGFkYXRhIjoie2ZvbzpiYXJ9Iiwicm9sZXMiOiIiLCJ1c2VyX2lkIjoiMTIzIn0.kL101I2ojpEvHEENfOUiREkiYNYUMP6x9xeAoeqgzhY"
		if token != expected {
			t.Fatal("expected token to not be empty", token)
		}
		if expires_at <= time.Date(2023, 1, 1, 0, 10, 0, 0, time.UTC).Unix() {
			t.Error("expected expires_at to be > 0")
		}
	})
	t.Run("should read claims from a valid access token", func(t *testing.T) {
		token, _, err := Generate(entities.CustomClaims{
			UserID:   "123",
			Roles:    []string{"some-role"},
			Metadata: "{foo:bar}",
		}, config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		claims, err := ReadClaims(token, config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if claims.UserID != "123" {
			t.Fatal("expected user_id to be 123", claims.UserID)
		}
		if claims.Metadata != "{foo:bar}" {
			t.Fatal("expected metadata to be {foo:bar}", claims.Metadata)
		}
		if claims.Roles[0] != "some-role" {
			t.Fatal("expected roles_values to be some-role", claims.Roles)
		}
	})
	t.Run("should return an error if the token is invalid", func(t *testing.T) {
		_, err := ReadClaims("token", config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("should return an error if the token is expired", func(t *testing.T) {
		token, _, err := Generate(entities.CustomClaims{
			UserID:   "123",
			Roles:    []string{""},
			Metadata: "{foo:bar}",
		}, config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		_, err = ReadClaims(token, config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err.Error() != apperrors.ErrAccessTokenExpired.Error() {
			t.Fatal("expected error to be 'access_token_expired'", err)
		}
	})
	t.Run("should return an error if the JWT secret is wrong", func(t *testing.T) {
		token, _, err := Generate(entities.CustomClaims{
			UserID:   "123",
			Roles:    []string{""},
			Metadata: "{foo:bar}",
		}, config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		_, err = ReadClaims(token, config.Config{
			JWT: config.JWTConfig{
				Secret:                     "xxxsecret2",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
