package service

import (
	"context"
	"fmt"

	"github.com/naavlad/auth-jwt-srv/internal/repository"
	"github.com/naavlad/auth-jwt-srv/internal/tokens"
)

type AuthService struct {
	repo         repository.Querier
	tokenManager *tokens.Manager
}

func NewAuthService(repo repository.Querier, tokenManager *tokens.Manager) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login проверяет учетные данные и возвращает токены
func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Проверка пароля (предполагаем, что в БД хранится открытый текст или хеш)
	if user.Password != password {
		return nil, fmt.Errorf("invalid credentials")
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken обновляет access токен
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Генерируем новый access токен
	accessToken, err := s.tokenManager.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

type UserInfo struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}

// GetUserInfo возвращает информацию о пользователе по токену
func (s *AuthService) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	claims, err := s.tokenManager.ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
