package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/herman-xphp/my-notes-api/internal/middleware"
	"github.com/herman-xphp/my-notes-api/internal/utils"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *fiber.App,
	authHandler *AuthHandler,
	noteHandler *NoteHandler,
	healthHandler *HealthHandler,
	jwtManager *utils.JWTManager,
) {
	// API v1 group
	api := app.Group("/api/v1")

	// Health check (no auth required)
	api.Get("/health", healthHandler.Check)

	// Auth routes (no auth required, but rate limited)
	auth := api.Group("/auth")
	auth.Post("/register", middleware.AuthRateLimit(), authHandler.Register)
	auth.Post("/login", middleware.AuthRateLimit(), authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected auth routes (require authentication)
	authProtected := auth.Group("", middleware.AuthMiddleware(jwtManager))
	authProtected.Post("/logout", authHandler.Logout)
	authProtected.Post("/logout-all", authHandler.LogoutAll)
	authProtected.Get("/me", authHandler.Me)

	// Note routes (all require authentication)
	notes := api.Group("/notes", middleware.AuthMiddleware(jwtManager))
	notes.Post("", noteHandler.Create)
	notes.Get("", noteHandler.GetAll)
	notes.Get("/:id", noteHandler.GetByID)
	notes.Put("/:id", noteHandler.Update)
	notes.Delete("/:id", noteHandler.Delete)
	notes.Post("/:id/restore", noteHandler.Restore)

	// Welcome route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to My Notes API",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Route not found",
		})
	})
}
