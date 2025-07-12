package integration

import "testing"

func TestProviderCallback(t *testing.T) {
	t.Run("unhappy scenarios: invalid data / generic", func(t *testing.T) {
		t.Run("calling GET /provider/callback returns an error if the provider is not enabled", func(t *testing.T) {})
		t.Run("calling GET /provider/callback with invalid state gets rejected", func(t *testing.T) {})
		t.Run("calling GET /provider/callback with invalid code gets rejected", func(t *testing.T) {})
	})
	t.Run("unhappy scenarios: error page", func(t *testing.T) {
		t.Run("calling GET /provider/callback returns to error page if user declines auth", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is using another method", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is blocked", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is deleted", func(t *testing.T) {})
		t.Run("calling GET /provider/callback returns to error page if user is not an early user", func(t *testing.T) {})
	})
	t.Run("happy scenarios", func(t *testing.T) {
		t.Run("calling GET /provider/callback gives [access_token, refresh_token] if the user already exists", func(t *testing.T) {})
		t.Run("calling GET /provider/callback gives [access_token, refresh_token] and creates user if the user does not exist", func(t *testing.T) {})
		t.Run("calling GET /provider/callback cleans the state", func(t *testing.T) {})
		t.Run("calling GET /provider/callback redirects to the welcome page", func(t *testing.T) {})
	})
}
