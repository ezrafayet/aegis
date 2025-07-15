package tokengen

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("should generate a random token with a prefix", func(t *testing.T) {
		token, err := Generate("test_", 8)
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
