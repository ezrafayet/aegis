package tokengen

import (
	"crypto/rand"
	"encoding/hex"
)

func Generate(prefix string, nPairs int) (string, error) {
	bytes := make([]byte, nPairs)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(bytes), nil
}
