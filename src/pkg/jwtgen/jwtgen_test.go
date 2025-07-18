package jwtgen

import (
	"aegis/pkg/apperrors"
	"testing"
	"time"
)

func TestAccessToken(t *testing.T) {
	t.Run("should create a new access token", func(t *testing.T) {
		token, expires_at, err := Generate(map[string]any{
			"user_id":         "123",
			"roles":           []string{""},
			"metadata_public": "{foo:bar}",
		}, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), 15, "app_name", "xxxsecret")
		if err != nil {
			t.Fatal(err)
		}
		expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhcHBfbmFtZSIsImV4cCI6MTY3MjUzMjEwMCwiaXNzIjoiYXBwX25hbWUiLCJpc3N1ZWRfYXQiOjE2NzI1MzEyMDAsIm1ldGFkYXRhX3B1YmxpYyI6Intmb286YmFyfSIsInJvbGVzIjpbIiJdLCJ1c2VyX2lkIjoiMTIzIn0.ur37iZg8pKHIHfsimsCiGbhzb6puq1HplJge3dWSW0A"
		if token != expected {
			t.Fatal("expected token to not be empty", token)
		}
		if expires_at <= time.Date(2023, 1, 1, 0, 10, 0, 0, time.UTC).Unix() {
			t.Error("expected expires_at to be > 0")
		}
	})
	t.Run("should read claims from a valid access token", func(t *testing.T) {
		token, _, err := Generate(map[string]any{
			"user_id":         "123",
			"roles":           []string{"some-role"},
			"metadata_public": "{foo:bar}",
		}, time.Now(), 15, "app_name", "xxxsecret")
		if err != nil {
			t.Fatal(err)
		}
		claims, err := ReadClaims(token, "xxxsecret")
		if err != nil {
			t.Fatal(err)
		}
		if claims["user_id"] != "123" {
			t.Fatal("expected user_id to be 123", claims["user_id"])
		}
		if claims["metadata_public"] != "{foo:bar}" {
			t.Fatal("expected metadata_public to be {foo:bar}", claims["metadata_public"])
		}
		roles := claims["roles"].([]interface{})
		if len(roles) == 0 || roles[0].(string) != "some-role" {
			t.Fatal("expected roles_values to be some-role", claims["roles"])
		}
	})
	t.Run("should return an error if the token is invalid", func(t *testing.T) {
		_, err := ReadClaims("token", "xxxsecret")
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("should return an error if the token is expired", func(t *testing.T) {
		token, _, err := Generate(map[string]any{
			"user_id":         "123",
			"roles":           []string{""},
			"metadata_public": "{foo:bar}",
		}, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), 15, "app_name", "xxxsecret")
		if err != nil {
			t.Fatal(err)
		}
		_, err = ReadClaims(token, "xxxsecret")
		if err.Error() != apperrors.ErrAccessTokenExpired.Error() {
			t.Fatal("expected error to be 'access_token_expired'", err)
		}
	})
	t.Run("should return an error if the JWT secret is wrong", func(t *testing.T) {
		token, _, err := Generate(map[string]any{
			"user_id":         "123",
			"roles":           []string{""},
			"metadata_public": "{foo:bar}",
		}, time.Now(), 15, "app_name", "xxxsecret")
		if err != nil {
			t.Fatal(err)
		}
		_, err = ReadClaims(token, "xxxsecret2")
		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
