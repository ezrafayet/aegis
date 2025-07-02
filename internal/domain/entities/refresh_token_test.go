package entities

import (
	"aegis/pkg/fingerprint"
	"testing"
	"time"
)

func TestRefreshToken(t *testing.T) {
	t.Run("should create a new refresh token", func(t *testing.T) {
		deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		token, _, _ := NewRefreshToken(User{ID: "123"}, deviceFingerprint, Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			}})
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
		if len(token.DeviceFingerprint) != 32 {
			t.Fatal("expected device_fingerprint to be 32 characters", token.DeviceFingerprint)
		}
	})
	t.Run("isExpired true", func(t *testing.T) {
		deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		token, _, _ := NewRefreshToken(User{ID: "123"}, deviceFingerprint, Config{
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
		token, _, _ := NewRefreshToken(User{ID: "123"}, "device-id", Config{
			JWT: JWTConfig{
				Secret:                     "xxxsecret",
				AccessTokenExpirationMin:   15,
				RefreshTokenExpirationDays: 30,
			}})
		if token.IsExpired() {
			t.Fatal("expected token to not be expired", token.ExpiresAt)
		}
	})
}

func TestGenerateDeviceFingerprint(t *testing.T) {
	t.Run("should generate a device fingerprint", func(t *testing.T) {
		deviceFingerprint, err := fingerprint.GenerateDeviceFingerprint("device-id")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		if deviceFingerprint == "" {
			t.Fatal("expected device fingerprint to not be empty", deviceFingerprint)
		}
	})
	t.Run("should generate a unique device fingerprint", func(t *testing.T) {
		deviceFingerprint1, err := fingerprint.GenerateDeviceFingerprint("device-id-1")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		deviceFingerprint2, err := fingerprint.GenerateDeviceFingerprint("device-id-2")
		if err != nil {
			t.Fatal("expected no error", err)
		}
		if deviceFingerprint1 == deviceFingerprint2 {
			t.Fatal("expected device fingerprint to be unique", deviceFingerprint1, deviceFingerprint2)
		}
	})
}
