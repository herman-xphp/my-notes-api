package impl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/dto"
	repomock "github.com/herman-xphp/my-notes-api/internal/repository/mock"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"github.com/herman-xphp/my-notes-api/internal/utils"
)

// ✅ FIX: Return service.AuthService interface
func setupAuthServiceTest() (service.AuthService, *repomock.UserRepository, *repomock.RefreshTokenRepository, *utils.JWTManager) {
	mockUserRepo := new(repomock.UserRepository)
	mockTokenRepo := new(repomock.RefreshTokenRepository)
	jwtManager := utils.NewJWTManager("test-secret-key-minimum-32-chars-long", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockUserRepo, mockTokenRepo, jwtManager)
	return svc, mockUserRepo, mockTokenRepo, jwtManager
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, mockUserRepo, mockTokenRepo, _ := setupAuthServiceTest()
	ctx := context.Background()

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "SecurePass123!",
	}

	// Mock expectations
	mockUserRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)
	mockUserRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	mockTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	// Execute
	result, err := svc.Register(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Equal(t, int64(900), result.ExpiresIn)
	assert.Equal(t, req.Name, result.User.Name)
	assert.Equal(t, req.Email, result.User.Email)
	assert.NotZero(t, result.User.ID)

	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	svc, mockUserRepo, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "existing@example.com",
		Password: "SecurePass123!",
	}

	// Mock expectations
	mockUserRepo.On("ExistsByEmail", ctx, req.Email).Return(true, nil)

	// Execute
	result, err := svc.Register(ctx, req)

	// Assert
	assert.Error(t, err)
	// ✅ FIX: Use service.ErrEmailAlreadyExists
	assert.Equal(t, service.ErrEmailAlreadyExists, err)
	assert.Nil(t, result)

	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "Create")
}

func TestAuthService_Register_WeakPassword(t *testing.T) {
	svc, mockUserRepo, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	req := dto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "weak", // Too weak
	}

	// Mock expectations
	mockUserRepo.On("ExistsByEmail", ctx, req.Email).Return(false, nil)

	// Execute
	result, err := svc.Register(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "password")

	mockUserRepo.AssertNotCalled(t, "Create")
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, mockUserRepo, mockTokenRepo, _ := setupAuthServiceTest()
	ctx := context.Background()

	password := "SecurePass123!"
	hashedPassword, _ := utils.HashPassword(password)

	existingUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: password,
	}

	// Mock expectations
	mockUserRepo.On("FindByEmail", ctx, req.Email).Return(existingUser, nil)
	mockTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	// Execute
	result, err := svc.Login(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Equal(t, existingUser.Name, result.User.Name)
	assert.Equal(t, existingUser.Email, result.User.Email)

	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc, mockUserRepo, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	req := dto.LoginRequest{
		Email:    "notfound@example.com",
		Password: "AnyPassword123!",
	}

	// Mock expectations
	mockUserRepo.On("FindByEmail", ctx, req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := svc.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidCredentials, err)
	assert.Nil(t, result)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc, mockUserRepo, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	correctPassword := "SecurePass123!"
	hashedPassword, _ := utils.HashPassword(correctPassword)

	existingUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	req := dto.LoginRequest{
		Email:    "john@example.com",
		Password: "WrongPassword123!",
	}

	// Mock expectations
	mockUserRepo.On("FindByEmail", ctx, req.Email).Return(existingUser, nil)

	// Execute
	result, err := svc.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidCredentials, err)
	assert.Nil(t, result)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	svc, mockUserRepo, mockTokenRepo, jwtManager := setupAuthServiceTest()
	ctx := context.Background()

	// Create a valid refresh token
	userID := uint(1)
	email := "john@example.com"
	refreshTokenString, _ := jwtManager.GenerateRefreshToken(userID, email)

	existingUser := &domain.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     email,
		CreatedAt: time.Now(),
	}

	storedToken := &domain.RefreshToken{
		ID:        1,
		UserID:    userID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	req := dto.RefreshTokenRequest{
		RefreshToken: refreshTokenString,
	}

	// Mock expectations
	mockTokenRepo.On("FindByToken", ctx, refreshTokenString).Return(storedToken, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	mockTokenRepo.On("Delete", ctx, refreshTokenString).Return(nil)
	mockTokenRepo.On("Create", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	// Execute
	result, err := svc.RefreshToken(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.NotEqual(t, refreshTokenString, result.RefreshToken) // Should be a new token
	assert.Equal(t, existingUser.Name, result.User.Name)

	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	svc, _, _, _ := setupAuthServiceTest()
	ctx := context.Background()

	req := dto.RefreshTokenRequest{
		RefreshToken: "invalid.token.string",
	}

	// Execute
	result, err := svc.RefreshToken(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidRefreshToken, err)
	assert.Nil(t, result)
}

func TestAuthService_RefreshToken_TokenNotInDatabase(t *testing.T) {
	svc, _, mockTokenRepo, jwtManager := setupAuthServiceTest()
	ctx := context.Background()

	// Create a valid JWT but not stored in database
	refreshTokenString, _ := jwtManager.GenerateRefreshToken(1, "john@example.com")

	req := dto.RefreshTokenRequest{
		RefreshToken: refreshTokenString,
	}

	// Mock expectations
	mockTokenRepo.On("FindByToken", ctx, refreshTokenString).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := svc.RefreshToken(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidRefreshToken, err)
	assert.Nil(t, result)

	mockTokenRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_ExpiredToken(t *testing.T) {
	svc, _, mockTokenRepo, jwtManager := setupAuthServiceTest()
	ctx := context.Background()

	// Create a refresh token
	refreshTokenString, _ := jwtManager.GenerateRefreshToken(1, "john@example.com")

	// Stored token is expired
	storedToken := &domain.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
		CreatedAt: time.Now(),
	}

	req := dto.RefreshTokenRequest{
		RefreshToken: refreshTokenString,
	}

	// Mock expectations
	mockTokenRepo.On("FindByToken", ctx, refreshTokenString).Return(storedToken, nil)
	mockTokenRepo.On("Delete", ctx, refreshTokenString).Return(nil)

	// Execute
	result, err := svc.RefreshToken(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidRefreshToken, err)
	assert.Nil(t, result)

	mockTokenRepo.AssertExpectations(t)
}

func TestAuthService_Logout_Success(t *testing.T) {
	svc, _, mockTokenRepo, _ := setupAuthServiceTest()
	ctx := context.Background()

	userID := uint(1)

	// Mock expectations
	mockTokenRepo.On("DeleteByUserID", ctx, userID).Return(nil)

	// Execute
	err := svc.Logout(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthService_LogoutAll_Success(t *testing.T) {
	svc, _, mockTokenRepo, _ := setupAuthServiceTest()
	ctx := context.Background()

	userID := uint(1)

	// Mock expectations
	mockTokenRepo.On("DeleteByUserID", ctx, userID).Return(nil)

	// Execute
	err := svc.LogoutAll(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockTokenRepo.AssertExpectations(t)
}
