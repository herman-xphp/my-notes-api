package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/herman-xphp/my-notes-api/internal/dto"
	"github.com/herman-xphp/my-notes-api/internal/middleware"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"github.com/herman-xphp/my-notes-api/internal/utils"
	"github.com/herman-xphp/my-notes-api/pkg/response"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService service.AuthService
	helper      *utils.HandlerHelper
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		helper:      utils.NewHandlerHelper(),
	}
}

// Register handler user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parse and validate
	if !h.helper.ParseAndValidate(c, &req) {
		return nil // Error already handled
	}
	// Call service
	result, err := h.authService.Register(c.Context(), req)
	if err != nil {
		return err // Will be handled by error middleware
	}
	return response.Created(c, "User registered successfully", result)
}

// Login handler user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if !h.helper.ParseAndValidate(c, &req) {
		return nil
	}

	// Call service
	result, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return err // Will be handled by error middleware
	}
	return response.Success(c, "Login successful", result)
}

// RefreshToken handler token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	if !h.helper.ParseAndValidate(c, &req) {
		return nil
	}
	// Call service
	result, err := h.authService.RefreshToken(c.Context(), req)
	if err != nil {
		return err // Will be handled by error middleware
	}
	return response.Success(c, "Token refreshed successfully", result)
}

// Logout handles user logout (single device)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Call service
	if err := h.authService.Logout(c.Context(), userID); err != nil {
		return err
	}

	return response.Success(c, "Logged out successfully", nil)
}

// LogoutAll handles user logout from all device
func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Call service
	if err := h.authService.LogoutAll(c.Context(), userID); err != nil {
		return err
	}

	return response.Success(c, "Logged out from all device successfully", nil)
}

// Me returns current user info
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	userEmail := middleware.GetUserEmail(c)

	// Return user info from token
	user := dto.UserResponse{
		ID:    userID,
		Email: userEmail,
	}

	return response.Success(c, "User retrieved successfully", user)
}
