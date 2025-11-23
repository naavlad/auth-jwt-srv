package handlers

import (
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

// Health godoc
// @Summary      Health check
// @Description  Проверка доступности сервиса
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}
