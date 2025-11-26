package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/utils"
)

type UserData struct {
	UserID       uuid.UUID `json:"id"`
	UserName     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

func (a *apiConfig) SignUp(w http.ResponseWriter, r *http.Request) {
	type SignUp struct {
		UserName        string `json:"username"`
		Email           string `json:"email"`
		Password        string `json:"password"`
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
		Email:    user_cred.Email,
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
		Email:    user_cred.Email,
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

	refreshToken := auth.MakeRefreshToken()
	go func() {
		err := a.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add((24 * time.Hour) * 30),
		})
		if err != nil {
			log.Printf("ERROR saving refresh token: %v", err)
		}
	}()

	resp := UserData{
		UserID:       user.ID,
		UserName:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		Role:         user.Role.String,
		CreatedAT:    user.CreatedAt,
		UpdatedAT:    user.UpdatedAt,
		AccessToken:  jwt,
		RefreshToken: refreshToken,
	}

	go func(name, email string) {
		SendWelcomeEmail(name, email)
	}(resp.UserName, resp.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(resp)
}

func (a *apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type loginCred struct {
		Email    string `json:"email"`
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	defer r.Body.Close()
	var user_cred loginCred
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user_cred); err != nil {
		utils.LogError(err.Error())
		utils.SendError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	data, err := a.db.GetUser(r.Context(), database.GetUserParams{
		Email:    user_cred.Email,
		Username: user_cred.UserName,
	})
	if err != nil {
		utils.LogError(err.Error())
		utils.SendError(w, "User does not exist", http.StatusNotFound)
		return
	}

	ok, err := utils.ComparePasswords(user_cred.Password, data.Password)
	if err != nil || !ok {
		if err != nil {
			go utils.LogError(err.Error())
		} else {
			go utils.LogError("password mismatch")
		}
		utils.SendError(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	jwt, _ := auth.MakeJWT(data.ID, a.jwtSecret)
	refreshToken := auth.MakeRefreshToken()
	go func() {
		err := a.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    data.ID,
			ExpiresAt: time.Now().Add((24 * time.Hour) * 30),
		})
		if err != nil {
			log.Printf("ERROR saving refresh token: %v", err)
		}
	}()

	resp := UserData{
		UserID:       data.ID,
		UserName:     data.Username,
		Email:        data.Email,
		Password:     data.Password,
		Role:         data.Role.String,
		UpdatedAT:    data.UpdatedAt,
		AccessToken:  jwt,
		RefreshToken: refreshToken,
	}

	go func(name, email string) {
		SendLoginWelcomeEmail(name, email)
	}(data.Username, data.Email)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
