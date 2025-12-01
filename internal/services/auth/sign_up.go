package authService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/database"
	utils2 "github.com/john-ayodeji/Linkrr/internal/utils"
)

var Cfg *config.ApiConfig

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

type signUp struct {
	UserName        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

var SignUpEvent = make(chan UserData, 100)

func SignUp(r *http.Request) (UserData, error, int) {
	var user_cred signUp
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user_cred); err != nil {
		return UserData{}, fmt.Errorf("something went wrong"), http.StatusInternalServerError
	}

	userEmail := strings.ToLower(user_cred.Email)
	u, _ := Cfg.Db.GetUser(r.Context(), database.GetUserParams{
		Email:    userEmail,
		Username: user_cred.UserName,
	})
	if u.ID != uuid.Nil {
		return UserData{}, fmt.Errorf("user already exists"), http.StatusUnauthorized
	}

	if user_cred.Password != user_cred.ConfirmPassword {
		return UserData{}, fmt.Errorf("passwords don't match"), http.StatusBadRequest
	}

	hashedPassword, err := utils2.HashPassword(user_cred.Password)
	if err != nil {
		return UserData{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}

	user, err := Cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
		Username: user_cred.UserName,
		Email:    userEmail,
		Password: hashedPassword,
	})
	if err != nil {
		go utils2.LogError(err.Error())
		return UserData{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}

	jwt, err := auth.MakeJWT(user.ID, Cfg.JWTSecret)
	if err != nil {
		go utils2.LogError(err.Error())
		return UserData{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}

	refreshToken := auth.MakeRefreshToken()
	go func() {
		err := Cfg.Db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
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

	SignUpEvent <- resp

	return resp, nil, 201
}
