package authService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/internal/database"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

var LoginEvent = make(chan UserData, 100)

func Login(r *http.Request) (UserData, error, int) {
	type loginCred struct {
		Email    string `json:"email"`
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	var user_cred loginCred
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user_cred); err != nil {
		return UserData{}, fmt.Errorf("something went wrong, try again later"), http.StatusInternalServerError
	}
	userEmail := strings.ToLower(user_cred.Email)

	data, err := Cfg.Db.GetUser(r.Context(), database.GetUserParams{
		Email:    userEmail,
		Username: user_cred.UserName,
	})
	if err != nil {
		return UserData{}, fmt.Errorf("user does not exist"), http.StatusNotFound
	}

	ok, err := utils.ComparePasswords(user_cred.Password, data.Password)
	if err != nil || !ok {
		return UserData{}, fmt.Errorf("incorrect password"), http.StatusUnauthorized
	}

	jwt, _ := auth.MakeJWT(data.ID, Cfg.JWTSecret)
	refreshToken := auth.MakeRefreshToken()
	go func() {
		err := Cfg.Db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
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

	LoginEvent <- resp

	return resp, nil, http.StatusOK
}
