package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(header http.Header) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	prefix := "Bearer"
	if !strings.HasPrefix(auth, prefix) {
		return "", fmt.Errorf("invalid authorization header")
	}

	token := strings.Fields(auth)[1]
	LogError(token)

	return token, nil
}
