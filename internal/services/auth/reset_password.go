package authService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

type request struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResetPasswordEmailData struct {
	Name  string
	Email string
}

var ResetPasswordEvent = make(chan ResetPasswordEmailData, 100)

func ResetPassword(r *http.Request) (Response, error, int) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return Response{}, fmt.Errorf("invalid password reset url"), http.StatusBadRequest
	}

	hashedToken := utils.HashToken(token)
	data, err := Cfg.Db.GetToken(r.Context(), hashedToken)
	if err != nil {
		return Response{}, fmt.Errorf("reset token does not exist"), http.StatusNotFound
	}
	if data.Used {
		return Response{}, fmt.Errorf("reset token has been used"), http.StatusUnauthorized
	}
	if data.ExpiresAt.Before(time.Now().UTC()) {
		return Response{}, fmt.Errorf("refresh token has expired"), http.StatusUnauthorized
	}

	var newCred request
	decoded := json.NewDecoder(r.Body)
	if err := decoded.Decode(&newCred); err != nil {
		return Response{}, fmt.Errorf("fill in all input fields"), http.StatusBadRequest
	}

	if newCred.Password != newCred.ConfirmPassword {
		return Response{}, fmt.Errorf("passwords don't match"), http.StatusBadRequest
	}

	hashedPassword, err := argon2id.CreateHash(newCred.Password, argon2id.DefaultParams)
	if err != nil {
		return Response{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}

	user, err := Cfg.Db.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Password: hashedPassword,
		ID:       data.UserID,
	})
	if err != nil {
		return Response{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}

	_ = Cfg.Db.SetUsed(r.Context())

	resp := Response{
		Status:  "success",
		Message: "Password has been changed successfully, check your email for a confirmation",
	}

	ResetPasswordEvent <- ResetPasswordEmailData{
		Name:  user.Username,
		Email: user.Email,
	}

	return resp, nil, http.StatusAccepted
}
