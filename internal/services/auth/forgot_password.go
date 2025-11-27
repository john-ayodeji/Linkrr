package authService

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/database"
	email2 "github.com/john-ayodeji/Linkrr/internal/services/email"
	"github.com/john-ayodeji/Linkrr/utils"
)

type cred struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func ForgotPassword(r *http.Request) (Response, error, int) {
	var userCred cred
	decoded := json.NewDecoder(r.Body)
	if err := decoded.Decode(&userCred); err != nil {
		return Response{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}
	userEmail := strings.ToLower(userCred.Email)
	user, err := Cfg.Db.GetUser(r.Context(), database.GetUserParams{
		Email:    userEmail,
		Username: userCred.Username,
	})
	if err != nil {
		return Response{}, fmt.Errorf("check your email for reset instructions if an account exists"), http.StatusNotFound
	}

	token := rand.Text()
	hashedToken := utils.HashToken(token)

	if err := Cfg.Db.CreateToken(r.Context(), database.CreateTokenParams{
		ID:          uuid.New(),
		UserID:      user.ID,
		HashedToken: hashedToken,
		ExpiresAt:   time.Now().UTC().Add(15 * time.Minute),
	}); err != nil {
		return Response{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}

	resetLink := fmt.Sprintf("%v/api/v1/auth/reset-password?token=%v", r.Host, token)
	go func(name, userEmail, url string) {
		email2.SendPasswordResetEmail(name, userEmail, url)
	}(user.Username, user.Email, resetLink)

	resp := Response{
		Status:  "success",
		Message: "Check your email for reset instructions if an account exists.",
	}

	return resp, nil, http.StatusOK
}
