package cookies

import (
	"aegix/internal/domain"
	"testing"
)

func TestNewCookie(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		cookie := newCookie("access_token", "test_value", 1717795200, domain.Config{})
		if cookie.Name != "access_token" {
			t.Errorf("expected cookie name to be 'access_token', got %s", cookie.Name)
		}
		if cookie.Value != "test_value" {
			t.Errorf("expected cookie value to be 'test_value', got %s", cookie.Value)
		}
		if cookie.Expires.Unix() != 1717795200 {
			t.Errorf("expected cookie expiration to be 1717795200, got %d", cookie.Expires.Unix())
		}
	})

	t.Run("with defaults 'true' ovverrides values", func(t *testing.T) {
		// todo
	})

	t.Run("with defaults 'false' does not ovverride values", func(t *testing.T) {
		// todo
	})
}
