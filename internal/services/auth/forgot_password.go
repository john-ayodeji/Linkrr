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
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

type cred struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ForgotPasswordEmailData struct {
	Name     string
	Email    string
	ResetURL string
}

var ForgotPasswordEvent = make(chan ForgotPasswordEmailData, 100)

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

	resp := Response{
		Status:  "success",
		Message: "Check your email for reset instructions if an account exists.",
	}

	ForgotPasswordEvent <- ForgotPasswordEmailData{
		Name:     user.Username,
		Email:    user.Email,
		ResetURL: resetLink,
	}

	return resp, nil, http.StatusOK
}
