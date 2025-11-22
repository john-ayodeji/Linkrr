package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/utils"
)

type UserData struct {
	UserID uuid.UUID `json:"id"`
	UserName string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (a *apiConfig) SignUp(w http.ResponseWriter, r *http.Request) {
	type SignUp struct {
		UserName string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	defer r.Body.Close()

	var user_cred SignUp
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user_cred); err != nil {
		go utils.LogError(err.Error())
		utils.SendError(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	u, _ := a.db.GetUser(r.Context(), database.GetUserParams{
		Email: user_cred.Email,
		Username: user_cred.UserName,
	})
	if u.ID != uuid.Nil {
		utils.SendError(w, "User already exists", http.StatusUnauthorized)
		return
	}

	if user_cred.Password != user_cred.ConfirmPassword {
		utils.SendError(w, "Passwords don't match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user_cred.Password)
	if err != nil {
		go utils.LogError(err.Error())
		utils.SendError(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	user, err := a.db.CreateUser(r.Context(), database.CreateUserParams{
		Username: user_cred.UserName,
		Email: user_cred.Email,
		Password: hashedPassword,
	})
	if err != nil {
		go utils.LogError(err.Error())
		return
	}

	jwt, err := auth.MakeJWT(user.ID, a.jwtSecret)
	if err != nil {
		go utils.LogError(err.Error())
		utils.SendError(w, "Something Went Wrong", http.StatusInternalServerError)
	}

	refresh_token := auth.MakeRefreshToken()

	resp := UserData {
		UserID: user.ID,
		UserName: user.Username,
		Email: user.Email,
		Password: user.Password,
		Role: user.Role.String,
		CreatedAT: user.CreatedAt,
		UpdatedAT: user.UpdatedAt,
		AccessToken: jwt,
		RefreshToken: refresh_token,
	}
	
	go func(name, email string) {
		SendWelcomeEmail(name, email)
	}(resp.UserName, resp.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(resp)
}