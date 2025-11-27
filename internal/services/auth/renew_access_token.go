package authService

import (
	"fmt"
	"net/http"
	"time"

	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/utils"
)

type Token struct {
	Token string `json:"token"`
}

func RenewAccessToken(r *http.Request) (Token, error, int) {
	token, err := utils.GetBearerToken(r.Header)
	if err != nil {
		return Token{}, err, http.StatusNotFound
	}

	data, err := Cfg.Db.VerifyRefreshToken(r.Context(), token)
	if err != nil {
		return Token{}, fmt.Errorf("refresh token has been revoked"), http.StatusUnauthorized
	}

	if data.ExpiresAt.Before(time.Now().UTC()) {
		return Token{}, fmt.Errorf("expired refresh token"), http.StatusUnauthorized
	}

	if data.RevokedAt.Valid {
		if data.RevokedAt.Time.Before(time.Now().UTC()) {
			return Token{}, fmt.Errorf("refresh token has been revoked"), http.StatusUnauthorized
		}
	}

	jwt, _ := auth.MakeJWT(data.UserID, Cfg.JWTSecret)
	return Token{Token: jwt}, nil, http.StatusAccepted
}
