package authHandler

import (
	"encoding/json"
	"net/http"

	authService "github.com/john-ayodeji/Linkrr/internal/services/auth"
	"github.com/john-ayodeji/Linkrr/utils"
)

func RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	resp, err, statusCode := authService.RenewAccessToken(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}
	r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	resp, err, statusCode := authService.RevokeRefreshToken(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
