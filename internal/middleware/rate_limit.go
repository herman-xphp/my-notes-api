package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitConfig holds rate limiter configuration
type RateLimitConfig struct {
	Max        int
	Expiration time.Duration
}

// RateLimit creates a rate limiter middleware
func RateLimit(config RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Rate limit by IP address
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

// AuthRateLimit creates stricter rate limiting for auth endpoints
func AuthRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:        5,                // 5 request
		Expiration: 15 * time.Minute, // per 15 minutes
	})
}

// GeneralRateLimit creates general rate limiting
func GeneralRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:        100,             // 100 request
		Expiration: 1 * time.Minute, // per minute
	})
}
