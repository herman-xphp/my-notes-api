package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/herman-xphp/my-notes-api/internal/domain"
	"github.com/herman-xphp/my-notes-api/internal/dto"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"github.com/herman-xphp/my-notes-api/internal/utils"
	"gorm.io/gorm"
)

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager       *utils.JWTManager
}

// Ensure authService implements service.AuthService
var _ service.AuthService = (*authService)(nil)

// NewAuthService creates a new auth service instance
func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *utils.JWTManager,
) service.AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, service.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Store refresh token
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Build response
	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresIn,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, service.ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Build response
	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresIn,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, service.ErrInvalidRefreshToken
	}

	// Find refresh token in database
	storedToken, err := s.refreshTokenRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}
	// Check if token expired
	if storedToken.IsExpired() {
		_ = s.refreshTokenRepo.Delete(ctx, req.RefreshToken)
		return nil, service.ErrInvalidRefreshToken
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Generate new token pair
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Delete old refresh token
	if err := s.refreshTokenRepo.Delete(ctx, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	// Store new refresh token
	newRefreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Build response
	return &dto.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresIn,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	return s.refreshTokenRepo.DeleteByUserID(ctx, userID)
}

func (s *authService) LogoutAll(ctx context.Context, userID uint) error {
	return s.refreshTokenRepo.DeleteByUserID(ctx, userID)
}
