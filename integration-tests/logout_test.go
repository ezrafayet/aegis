package integration

import "testing"

func TestLogout(t *testing.T) {
	t.Run("calling GET /logout sets zero cookies", func(t *testing.T) {})
	t.Run("calling GET /logout with a refresh token deletes the refresh token", func(t *testing.T) {})
	t.Run("calling GET /logout without a refresh token does not break", func(t *testing.T) {})
}
