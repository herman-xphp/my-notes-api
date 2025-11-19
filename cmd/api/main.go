package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/herman-xphp/my-notes-api/configs"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	"github.com/herman-xphp/my-notes-api/internal/repository/mysql"
	"github.com/herman-xphp/my-notes-api/pkg/database"
	"github.com/herman-xphp/my-notes-api/pkg/response"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting %s in %s mode...", cfg.App.Name, cfg.App.Env)

	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDatabase(db)

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	initRepositories(db)
	log.Println("Repositories initialized")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: errorHandler,
	})

	// Setup middleware
	setupMiddlewares(app, cfg)

	// Setup routes
	setupRoutes(app)

	// Start server
	go startServer(app, cfg)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func connectDatabase(cfg *configs.Config) (*gorm.DB, error) {
	logLevel := gormlogger.Info
	if cfg.IsProduction() {
		logLevel = gormlogger.Error
	}

	db, err := database.NewMySQL(database.Config{
		DSN:          cfg.GetDBDSN(),
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogLevel:     logLevel,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func closeDatabase(db *gorm.DB) {
	if err := database.Close(db); err != nil {
		log.Printf("Failed to close database: %v", err)
	}
}

func setupMiddlewares(app *fiber.App, cfg *configs.Config) {
	// Recover from panics
	app.Use(recover.New())

	// Logger middleware
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// CROS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
}

func setupRoutes(app *fiber.App) {
	// API v1 routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, "Server is healthy", fiber.Map{
			"status": "ok",
			"env":    os.Getenv("APP_ENV"),
		})
	})

	// Welcome route
	app.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, "Welcome to My Notes API", fiber.Map{
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
		})
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return response.NotFound(c, "Route not found")
	})
}

func startServer(app *fiber.App, cfg *configs.Config) {
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("Server is running on http://localhost%s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return response.Error(c, code, "An error occurred", err.Error())
}

func initRepositories(db *gorm.DB) *repository.Repositories {
	return &repository.Repositories{
		User:         mysql.NewUserRepository(db),
		Note:         mysql.NewNoteRepository(db),
		RefreshToken: mysql.NewRefreshTokenRepository(db),
	}
}
