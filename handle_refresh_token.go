package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/john-ayodeji/Linkrr/internal/auth"
	"github.com/john-ayodeji/Linkrr/utils"
)

func (a *apiConfig) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	defer r.Body.Close()
	token, err := utils.GetBearerToken(r.Header)
	if err != nil {
		utils.LogError(err.Error())
		utils.SendError(w, err.Error(), http.StatusNotFound)
		return
	}

	data, err := a.db.VerifyRefreshToken(r.Context(), token)
	if err != nil {
		utils.SendError(w, "Refresh token has been revoked", http.StatusNotFound)
		return
	}

	if data.ExpiresAt.Before(time.Now().UTC()) {
		utils.SendError(w, "Expired refresh token", http.StatusUnauthorized)
		return
	}

	if data.RevokedAt.Valid {
		if data.RevokedAt.Time.Before(time.Now().UTC()) {
			utils.SendError(w, "Refresh token has been revoked", http.StatusUnauthorized)
			return
		}
	}

	jwt, _ := auth.MakeJWT(data.UserID, a.jwtSecret)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{Token: jwt})
}

func (a *apiConfig) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}

	defer r.Body.Close()
	token, err := utils.GetBearerToken(r.Header)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusNotFound)
		return
	}

	_, err1 := a.db.VerifyRefreshToken(r.Context(), token)
	if err1 != nil {
		utils.SendError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	Err := a.db.RevokeRefreshToken(r.Context(), token)
	if Err != nil {
		utils.SendError(w, "Something wen't wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response{Status: "success"})
}
