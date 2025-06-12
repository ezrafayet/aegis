package domain

import (
	"strings"
	"testing"
)

func TestGenerateRandomToken(t *testing.T) {
	t.Run("should generate a random token with a prefix", func(t *testing.T) {
		token, err := GenerateRandomToken("test_", 8)
		if err != nil {
			t.Fatal(err)
		}
		if len(token) != 21 {
			t.Fatal("expected token to have length 21", len(token))
		}
		if !strings.HasPrefix(token, "test_") {
			t.Fatal("expected token to have prefix test", token)
		}
	})
}
