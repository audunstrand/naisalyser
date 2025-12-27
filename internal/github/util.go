package github

import (
	"encoding/base64"
	"strings"
)

// decodeBase64 decodes a base64 encoded string (handles line breaks)
func decodeBase64(encoded string) (string, error) {
	// GitHub returns base64 with newlines
	cleaned := strings.ReplaceAll(encoded, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\"", "")

	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
