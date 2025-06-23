package cookies

import (
	"othnx/internal/infrastructure/config"
	"testing"
)

func TestNewCookie(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		cookie := newCookie("access_token", "test_value", 1717795200, config.Config{})
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

	t.Run("is zero cookie detects a zero cookie", func(t *testing.T) {
		cookie := NewAccessCookieZero(config.Config{})
		if !IsZeroCookie(cookie) {
			t.Errorf("expected cookie to be zero")
		}
	})

	t.Run("is zero cookie detects a non-zero cookie", func(t *testing.T) {
		cookie := NewAccessCookie("test_value", 1717795200, config.Config{})
		if IsZeroCookie(cookie) {
			t.Errorf("expected cookie to be non-zero")
		}
	})
}
