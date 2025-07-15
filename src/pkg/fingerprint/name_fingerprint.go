package fingerprint

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GenerateNameFingerprint(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", errors.New("empty_name")
	}
	normalized := strings.ToLower(trimmed)
	transformer := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, err := transform.String(transformer, normalized)
	if err != nil {
		result = normalized
	}
	result = strings.Join(strings.Fields(result), " ")
	hash := md5.Sum([]byte(result))
	return hex.EncodeToString(hash[:]), nil
}
