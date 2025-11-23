package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"secure_password"`
}

// Login godoc
// @Summary      Аутентификация пользователя
// @Description  Получение JWT токенов по username и password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Учетные данные"
// @Success      200      {object}  service.LoginResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	resp, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// Refresh godoc
// @Summary      Обновление access токена
// @Description  Получение нового access токена с помощью refresh токена
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Refresh токен"
// @Success      200      {object}  RefreshResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	accessToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	writeJSON(w, http.StatusOK, RefreshResponse{AccessToken: accessToken})
}

// Me godoc
// @Summary      Получение информации о пользователе
// @Description  Возвращает информацию о текущем пользователе по access токену
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  service.UserInfo
// @Failure      401  {object}  ErrorResponse
// @Router       /auth/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, "authorization header required")
		return
	}

	// Ожидаем формат "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		writeError(w, http.StatusUnauthorized, "invalid authorization header format")
		return
	}

	accessToken := parts[1]

	userInfo, err := h.authService.GetUserInfo(r.Context(), accessToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid access token")
		return
	}

	writeJSON(w, http.StatusOK, userInfo)
}
