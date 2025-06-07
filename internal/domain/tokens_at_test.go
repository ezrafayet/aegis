package domain

import (
	"testing"
	"time"
)

func TestAccessToken(t *testing.T) {
	t.Run("should create a new access token", func(t *testing.T) {
		token, expires_at, err := NewAccessToken(CustomClaims{
			UserID:   "123",
			Roles:    []string{""},
			Metadata: "{foo:bar}",
		}, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJleHAiOjE2NzI1MzIxMDAsImlzcyI6IiIsImlzc3VlZF9hdCI6MTY3MjUzMTIwMCwibWV0YWRhdGEiOiJ7Zm9vOmJhcn0iLCJ1c2VyX2lkIjoiMTIzIn0.OVxAoDWCBvPbCZOTfV7ReSOAlGnV-gicgInpbKTihA0"
		if token != expected {
			t.Fatal("expected token to not be empty", token)
		}
		if expires_at <= time.Date(2023, 1, 1, 0, 10, 0, 0, time.UTC).Unix() {
			t.Error("expected expires_at to be > 0")
		}
	})
	t.Run("should read claims from a valid access token", func(t *testing.T) {
		token, _, err := NewAccessToken(CustomClaims{
			UserID:   "123",
			Roles:    []string{""},
			Metadata: "{foo:bar}",
		}, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		claims, err := ReadAccessTokenClaims(token, Config{
			JWT: JWTConfig{
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
	})
	t.Run("should return an error if the token is invalid (not even a jwt)", func(t *testing.T) {
		_, err := ReadAccessTokenClaims("token", Config{
			JWT: JWTConfig{
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
		token, _, err := NewAccessToken(CustomClaims{
			UserID:   "123",
			Roles:    []string{""},
			Metadata: "{foo:bar}",
		}, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		_, err = ReadAccessTokenClaims(token, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if err.Error() != "access_token_expired" {
			t.Fatal("expected error to be 'access_token_expired'", err)
		}
	})
	t.Run("should return an error if the JWT secret is different", func(t *testing.T) {
		token, _, err := NewAccessToken(CustomClaims{
			UserID:   "123",
			Roles:    []string{""},
			Metadata: "{foo:bar}",
		}, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		}, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		_, err = ReadAccessTokenClaims(token, Config{
			JWT: JWTConfig{
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
