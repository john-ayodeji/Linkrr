package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "linkrr",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   userID.String(),
	})
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ValidateJWT(tokenStr, tokenSecret string) (*jwt.Token, *jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			// Ensure HS256 is used
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
			}
			return []byte(tokenSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithLeeway(0),
	)
	if err != nil {
		return nil, nil, err
	}
	if !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}
	return token, claims, nil
}

func ValidateJWTHelper(tokenStr, jwtSecret string) (bool, *jwt.RegisteredClaims, error) {
	token, claims, err := ValidateJWT(tokenStr, jwtSecret)
	if err != nil {
		return false, nil, fmt.Errorf("%v", err)
	}
	if !token.Valid {
		return false, nil, fmt.Errorf("invalid authorization token")
	}

	if claims.ExpiresAt.Before(time.Now().UTC()) {
		return false, nil, fmt.Errorf("expired authorization token")
	}

	return true, claims, nil
}
