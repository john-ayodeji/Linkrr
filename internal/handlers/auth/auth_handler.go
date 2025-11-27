package authHandler

import (
	"encoding/json"
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/config"
	"github.com/john-ayodeji/Linkrr/internal/services/auth"
	"github.com/john-ayodeji/Linkrr/utils"
)

var Cfg *config.ApiConfig

func SignUp(w http.ResponseWriter, r *http.Request) {
	resp, err, statusCode := authService.SignUp(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	_ = json.NewEncoder(w).Encode(resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	resp, err, statusCode := authService.Login(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
