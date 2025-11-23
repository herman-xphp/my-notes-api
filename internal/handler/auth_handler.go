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
	validator   *utils.Validator
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   utils.NewValidator(),
	}
}

// Register handler user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
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

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
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

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
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
