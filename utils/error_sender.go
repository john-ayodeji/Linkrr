package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func SendError(w http.ResponseWriter, errMessage string, statusCode int) {
	errData := Error {
		Error: errMessage,
	}

	data, err := json.Marshal(errData)
	if err != nil {
		log.Printf("Error encoding error response: %s", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(Error { Error: "Something went wrong" })
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(data)
}