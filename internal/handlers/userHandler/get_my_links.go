package userHandler

import (
	"encoding/json"
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/services/users"
	"github.com/john-ayodeji/Linkrr/internal/utils"
)

func GetMyLinks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp, err, statusCode := users.GetMyLinks(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
