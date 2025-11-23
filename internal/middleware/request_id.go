package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get response ID from header or generate new one
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context and response header
		c.Locals("RequestID", requestID)
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Locals("RequestID").(string)
	if !ok {
		return ""
	}

	return requestID
}
