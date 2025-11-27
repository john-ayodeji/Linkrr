package authHandler

import (
	"encoding/json"
	"net/http"

	authService "github.com/john-ayodeji/Linkrr/internal/services/auth"
	"github.com/john-ayodeji/Linkrr/utils"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	resp, err, statusCode := authService.ForgotPassword(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	resp, err, statusCode := authService.ResetPassword(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
