package shortenerHandler

import (
	"encoding/json"
	"net/http"

	"github.com/john-ayodeji/Linkrr/internal/services/shortener"
	"github.com/john-ayodeji/Linkrr/utils"
)

func CreateAlias(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resp, err, statusCode := shortener.CreateAlias(r)
	if err != nil {
		utils.SendError(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
