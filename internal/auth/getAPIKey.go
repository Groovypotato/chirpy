package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHdr := headers.Get("Authorization")
	if authHdr == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.Fields(authHdr)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "ApiKey") {
		return "", errors.New("invalid authorization header format")
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("empty API key")
	}
	return token, nil
}