package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// Health godoc
// @Summary Health check
// @Description Returns the health status of the application
// @Tags info
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	}
}
