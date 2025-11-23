package tokens

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	manager := NewManager("test-secret-key", 15*time.Minute, 168*time.Hour)

	token, err := manager.GenerateAccessToken(1, "testuser")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := NewManager("test-secret-key", 15*time.Minute, 168*time.Hour)

	token, err := manager.GenerateRefreshToken(1, "testuser")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	manager := NewManager("test-secret-key", 15*time.Minute, 168*time.Hour)

	token, err := manager.GenerateAccessToken(1, "testuser")
	assert.NoError(t, err)

	claims, err := manager.ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, int32(1), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	manager := NewManager("test-secret-key", 15*time.Minute, 168*time.Hour)

	claims, err := manager.ValidateToken("invalid-token")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	manager := NewManager("test-secret-key", -1*time.Hour, 168*time.Hour)

	token, err := manager.GenerateAccessToken(1, "testuser")
	assert.NoError(t, err)

	claims, err := manager.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}
