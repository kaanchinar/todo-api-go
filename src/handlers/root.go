package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kaanchinar/todo-app/config"
)

type RootResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Root godoc
// @Summary API info
// @Description Get the application name and version
// @Tags info
// @Produce json
// @Success 200 {object} RootResponse
// @Router / [get]

func Root(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RootResponse{
			Name:    cfg.AppName,
			Version: cfg.AppVersion,
		})
	}
}
