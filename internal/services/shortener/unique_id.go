package shortener

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

func GenerateRandomURLID(n int) string {
	const chars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	b := make([]byte, n)
	_, _ = rand.Read(b)

	for i := 0; i < 6; i++ {
		b[i] = chars[b[i]%byte(len(chars))]
	}

	return string(b)
}

func ShortURLExists(r *http.Request, urlCode string) bool {
	data, err := Cfg.Db.GetURL(r.Context(), urlCode)
	if err != nil {
		return false
	}
	if data.ShortCode == urlCode {
		return true
	}
	return false
}

func RandomIDwithRetry(r *http.Request) (string, error) {
	maxAttempt := 5
	for i := 0; i <= maxAttempt; i++ {
		code := GenerateRandomURLID(6)
		exists := ShortURLExists(r, code)
		if !exists {
			return code, nil
		}
	}
	return "", fmt.Errorf("something went wrong, try again later")
}
