package utils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func SendError(w http.ResponseWriter, errMessage string, statusCode int) {
	errData := Error{
		Error: errMessage,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(errData)
}
