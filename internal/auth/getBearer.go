package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHdr := headers.Get("Authorization")
	if authHdr == "" {
		return "", errors.New("authorization header missing")
	}
	parts := strings.Fields(authHdr)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header format")
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}
