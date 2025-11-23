package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"github.com/herman-xphp/my-notes-api/pkg/response"
	"github.com/rs/zerolog/log"
)

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	// Fiber error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	// Business logic errors
	switch {
	case errors.Is(err, service.ErrEmailAlreadyExists):
		code = fiber.StatusConflict
		message = "Email already exists"
	case errors.Is(err, service.ErrInvalidCredentials):
		code = fiber.StatusUnauthorized
		message = "Invalid email or password"
	case errors.Is(err, service.ErrUserNotFound):
		code = fiber.StatusNotFound
		message = "User not found"
	case errors.Is(err, service.ErrNoteNotFound):
		code = fiber.StatusNotFound
		message = "Note not found"
	case errors.Is(err, service.ErrInvalidRefreshToken):
		code = fiber.StatusUnauthorized
		message = "Invalid or expired refresh token"
	case errors.Is(err, service.ErrUnauthorizedAccess):
		code = fiber.StatusForbidden
		message = "Unauthorized access to resource"
	}

	// Log error
	if code >= 500 {
		log.Error().
			Err(err).
			Str("path", c.Path()).
			Str("method", c.Method()).
			Int("status", code).
			Msg("Internal server error")
	}
	// Return error response
	return response.Error(c, code, message, nil)
}
