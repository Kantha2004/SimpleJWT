package api

import (
	"encoding/json"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "ok",
		"service": "SimpleJWT",
		"message": "Application is running...",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
