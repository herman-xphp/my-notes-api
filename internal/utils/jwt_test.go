package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTManager(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)
	assert.NotNil(t, manager)
}

func TestGenerateAccessToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	token, err := manager.GenerateAccessToken(1, "test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	token, err := manager.GenerateRefreshToken(1, "test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateTokenPair(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	pair, err := manager.GenerateTokenPair(1, "test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, pair)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Equal(t, int64(900), pair.ExpiresIn) // 15 minutes
}

func TestValidateToken_Success(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	// Generate token
	token, err := manager.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	// Validate token
	claims, err := manager.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	_, err := manager.ValidateToken("invalid.token.string")
	assert.Error(t, err)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	manager1 := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)
	manager2 := NewJWTManager("different-secret-key-minimum-32", 15*time.Minute, 7*24*time.Hour)

	// Generate with manager1
	token, err := manager1.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	// Try to validate with manager2 (different secret)
	_, err = manager2.ValidateToken(token)
	assert.Error(t, err)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Create manager with very short expiry
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 1*time.Millisecond, 7*24*time.Hour)

	// Generate token
	token, err := manager.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to validate expired token
	_, err = manager.ValidateToken(token)
	assert.Error(t, err)
}

func TestExtractUserID(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	// Generate token
	token, err := manager.GenerateAccessToken(123, "test@example.com")
	require.NoError(t, err)

	// Extract user ID
	userID, err := manager.ExtractUserID(token)
	assert.NoError(t, err)
	assert.Equal(t, uint(123), userID)
}

func TestIsTokenExpired(t *testing.T) {
	// Test valid token
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)
	token, err := manager.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	expired, err := manager.IsTokenExpired(token)
	assert.NoError(t, err)
	assert.False(t, expired)

	// Test expired token
	shortManager := NewJWTManager("test-secret-key-minimum-32-chars", 1*time.Millisecond, 7*24*time.Hour)
	expiredToken, err := shortManager.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	expired, err = shortManager.IsTokenExpired(expiredToken)
	// Token validation will fail for expired token
	assert.Error(t, err)
}

func TestGetTokenRemainingTime(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	token, err := manager.GenerateAccessToken(1, "test@example.com")
	require.NoError(t, err)

	remaining, err := manager.GetTokenRemainingTime(token)
	assert.NoError(t, err)
	assert.Greater(t, remaining, 14*time.Minute) // Should be close to 15 minutes
	assert.LessOrEqual(t, remaining, 15*time.Minute)
}

func TestTokenClaimsFields(t *testing.T) {
	manager := NewJWTManager("test-secret-key-minimum-32-chars", 15*time.Minute, 7*24*time.Hour)

	token, err := manager.GenerateAccessToken(42, "user@example.com")
	require.NoError(t, err)

	claims, err := manager.ValidateToken(token)
	require.NoError(t, err)

	// Check all claims
	assert.Equal(t, uint(42), claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
	assert.NotEmpty(t, claims.ID) // JWT ID should be set
	assert.Equal(t, "my-notes-api", claims.Issuer)
	assert.Equal(t, "42", claims.Subject)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.NotBefore)
}
