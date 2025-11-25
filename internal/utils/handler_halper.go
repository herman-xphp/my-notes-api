package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/herman-xphp/my-notes-api/pkg/response"
)

// HandlerHelper provides reusable handler utilities
type HandlerHelper struct {
	validator *Validator
}

// NewHandlerHelper creates a new handler helper instance
func NewHandlerHelper() *HandlerHelper {
	return &HandlerHelper{
		validator: NewValidator(),
	}
}

// ParseAndValidate parses request body and validates it
// Returns true if successful, false if error (and already sent response)
func (h *HandlerHelper) ParseAndValidate(c *fiber.Ctx, req interface{}) bool {
	// Parse request body
	if err := c.BodyParser(req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return false
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return false
	}

	return true
}

// ParseQuery parses query parameters and validates
func (h *HandlerHelper) ParseQuery(c *fiber.Ctx, req interface{}) bool {
	// Parse query params
	if err := c.QueryParser(req); err != nil {
		response.BadRequest(c, "Invalid query parameters", err.Error())
		return false
	}

	// Validate request
	if err := h.validator.Validate(req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return false
	}

	return true
}

// ParseID parses ID from URL parameter
func (h *HandlerHelper) ParseID(c *fiber.Ctx, paramName string) (uint, bool) {
	idStr := c.Params(paramName)
	if idStr == "" {
		response.BadRequest(c, "Missing ID parameter", nil)
		return 0, false
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID format", err.Error())
		return 0, false
	}

	return uint(id), true
}

// GetUserID extracts user ID from context (from auth middleware)
func (h *HandlerHelper) GetUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}

// GetUserEmail extracts user email from context
func (h *HandlerHelper) GetUserEmail(c *fiber.Ctx) string {
	email, ok := c.Locals("userEmail").(string)
	if !ok {
		return ""
	}
	return email
}

// GetRequestID extracts request ID from context
func (h *HandlerHelper) GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Locals("requestID").(string)
	if !ok {
		return ""
	}
	return requestID
}

// MustGetUserID panics if user ID not found (for protected routes only)
func (h *HandlerHelper) MustGetUserID(c *fiber.Ctx) uint {
	userID := h.GetUserID(c)
	if userID == 0 {
		panic("user ID not found in context - ensure auth middleware is applied")
	}
	return userID
}
