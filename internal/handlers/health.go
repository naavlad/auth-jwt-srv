package handlers

import (
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// Health проверяет доступность сервиса
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}
