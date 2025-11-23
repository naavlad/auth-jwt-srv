package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Устанавливаем переменные окружения для теста
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("JWT_ACCESS_TOKEN_DURATION", "30m")
	os.Setenv("JWT_REFRESH_TOKEN_DURATION", "240h")
	os.Setenv("SERVER_PORT", "8090")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_ACCESS_TOKEN_DURATION")
		os.Unsetenv("JWT_REFRESH_TOKEN_DURATION")
		os.Unsetenv("SERVER_PORT")
	}()

	cfg, err := Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres://test:test@localhost:5432/testdb", cfg.Database.URL)
	assert.Equal(t, "test-secret-key", cfg.JWT.Secret)
	assert.Equal(t, 30*time.Minute, cfg.JWT.AccessTokenDuration)
	assert.Equal(t, 240*time.Hour, cfg.JWT.RefreshTokenDuration)
	assert.Equal(t, "8090", cfg.Server.Port)
}

func TestLoad_MissingRequired(t *testing.T) {
	// Убираем обязательные переменные
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")

	cfg, err := Load()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoad_DefaultValues(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, err := Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Server.Port) // default value
}
