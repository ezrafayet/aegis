package domain

import (
	"testing"
	"time"
)

func TestRefreshToken(t *testing.T) {
	t.Run("should create a new refresh token", func(t *testing.T) {
		token, _ := NewRefreshToken(User{ID: "123"}, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if token.UserID != "123" {
			t.Fatal("expected user_id to be 123", token.UserID)
		}
		if token.Token == "" {
			t.Fatal("expected token to not be empty", token.Token)
		}
		if token.CreatedAt.IsZero() {
			t.Fatal("expected created_at to not be zero", token.CreatedAt)
		}
		if token.ExpiresAt.IsZero() {
			t.Fatal("expected expires_at to not be zero", token.ExpiresAt)
		}
		if token.ExpiresAt.Before(time.Now()) {
			t.Fatal("expected expires_at to be in the future", token.ExpiresAt)
		}
	})
	t.Run("isExpired true", func(t *testing.T) {
		token, _ := NewRefreshToken(User{ID: "123"}, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		token.ExpiresAt = time.Now().Add(-time.Hour * 24)
		if !token.IsExpired() {
			t.Fatal("expected token to be expired", token.ExpiresAt)
		}
	})
	t.Run("isExpired false", func(t *testing.T) {
		token, _ := NewRefreshToken(User{ID: "123"}, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			},
		})
		if token.IsExpired() {
			t.Fatal("expected token to not be expired", token.ExpiresAt)
		}
	})
}
