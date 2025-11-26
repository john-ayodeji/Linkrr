package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func MakeRefreshToken() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("REFRESH TOKEN ERROR: %s", err)
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	return token[:12]
}
