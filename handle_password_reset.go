package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/utils"
)

func (a *apiConfig) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	type response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	var userCred request
	defer r.Body.Close()
	decoded := json.NewDecoder(r.Body)
	if err := decoded.Decode(&userCred); err != nil {
		utils.SendError(w, "Something went wrong, try again later", http.StatusInternalServerError)
		return
	}

	user, err := a.db.GetUser(r.Context(), database.GetUserParams{
		Email:    userCred.Email,
		Username: userCred.Username,
	})
	if err != nil {
		utils.SendError(w, "User with email/username does not exist", http.StatusNotFound)
		return
	}

	token := rand.Text()
	hashedToken := utils.HashToken(token)

	if err := a.db.CreateToken(r.Context(), database.CreateTokenParams{
		ID:          uuid.New(),
		UserID:      user.ID,
		HashedToken: hashedToken,
		ExpiresAt:   time.Now().UTC().Add(15 * time.Minute),
	}); err != nil {
		utils.SendError(w, "Something went wrong, try again later", http.StatusInternalServerError)
		return
	}

	resetLink := fmt.Sprintf("%v/reset-password?token=%v", a.baseUrl, token)
	go func(name, email, url string) {
		SendPasswordResetEmail(name, email, url)
	}(user.Username, user.Email, resetLink)

	resp := response{
		Status:  "success",
		Message: "Check your email for reset instructions if an account exists.",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (a *apiConfig) ResetPassword(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	type response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	defer r.Body.Close()

	token := r.URL.Query().Get("token")
	if token == "" {
		utils.SendError(w, "Invalid password reset url", http.StatusBadRequest)
		return
	}

	hashedToken := utils.HashToken(token)
	data, err := a.db.GetToken(r.Context(), hashedToken)
	if err != nil {
		utils.SendError(w, "reset token does not exist", http.StatusNotFound)
		return
	}
	if data.Used {
		utils.SendError(w, "reset token has been used", http.StatusUnauthorized)
		return
	}
	if data.ExpiresAt.Before(time.Now().UTC()) {
		utils.SendError(w, "reset token has expired", http.StatusUnauthorized)
		return
	}

	var newCred request
	decoded := json.NewDecoder(r.Body)
	if err := decoded.Decode(&newCred); err != nil {
		utils.SendError(w, "Fill in all input fields", http.StatusInternalServerError)
		return
	}

	if newCred.Password != newCred.ConfirmPassword {
		utils.SendError(w, "Passwords don't match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := argon2id.CreateHash(newCred.Password, argon2id.DefaultParams)
	if err != nil {
		utils.SendError(w, "Something went wrong, try again later", http.StatusInternalServerError)
		return
	}

	user, err := a.db.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Password: hashedPassword,
		ID:       data.UserID,
	})
	if err != nil {
		utils.SendError(w, "Something went wrong, try again later", http.StatusInternalServerError)
		return
	}

	_ = a.db.SetUsed(r.Context())

	go func(name, email string) {
		SendPasswordChangedEmail(name, email)
	}(user.Username, user.Email)

	resp := response{
		Status:  "success",
		Message: "Password has been changed successfully, check your email for a confirmation",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
