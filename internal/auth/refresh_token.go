package auth

import (
	"crypto/rand"
)

func MakeRefreshToken() string {
	token := rand.Text()
	return token
}
