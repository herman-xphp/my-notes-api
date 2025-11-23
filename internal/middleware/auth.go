package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/herman-xphp/my-notes-api/internal/utils"
	"github.com/herman-xphp/my-notes-api/pkg/response"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtManager *utils.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authentication header
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		// Check Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization header format")
		}

		token := parts[1]

		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired token")
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *fiber.Ctx) string {
	email, ok := c.Locals("userEmail").(string)
	if !ok {
		return ""
	}
	return email
}
