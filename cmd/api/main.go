package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	"github.com/herman-xphp/my-notes-api/configs"
	"github.com/herman-xphp/my-notes-api/internal/handler"
	"github.com/herman-xphp/my-notes-api/internal/middleware"
	"github.com/herman-xphp/my-notes-api/internal/repository"
	"github.com/herman-xphp/my-notes-api/internal/repository/mysql"
	"github.com/herman-xphp/my-notes-api/internal/service"
	"github.com/herman-xphp/my-notes-api/internal/service/impl"
	"github.com/herman-xphp/my-notes-api/internal/utils"
	"github.com/herman-xphp/my-notes-api/pkg/database"
)

func main() {
	// Load configuration dengan validasi
	cfg, err := configs.LoadConfig()
	if err != nil {
		// Log error dengan detail yang jelas
		log.Fatalf("❌ Failed to load configuration: %v\n\n"+
			"💡 Tips:\n"+
			"   1. Copy .env.example ke .env: cp .env.example .env\n"+
			"   2. Generate JWT_SECRET: go run scripts/generate-secret.go\n"+
			"   3. Update .env dengan JWT_SECRET yang baru\n"+
			"   4. Set environment variables lainnya sesuai kebutuhan\n",
			err)
	}

	// Print startup info
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  🚀 Starting %s\n", cfg.App.Name)
	fmt.Printf("║  📦 Environment: %s\n", cfg.App.Env)
	fmt.Printf("║  🔌 Port: %s\n", cfg.App.Port)
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	fmt.Println("✅ Database connected successfully")

	// Run migrations (implement ini sesuai kebutuhan)
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	fmt.Println("✅ Database migrations completed")

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	// Initialize repositories
	repos := initRepositories(db)
	log.Println("✅ Repositories initialized")

	// Initialize services
	services := initServices(repos, jwtManager)
	log.Println("✅ Services initialized")

	// Initialize handlers
	handlers := initHandlers(services, db)
	log.Println("✅ Handlers initialized")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: middleware.ErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Setup middlewares
	setupMiddlewares(app, cfg)

	// Setup routes
	handler.SetupRoutes(
		app,
		handlers.auth,
		handlers.note,
		handlers.health,
		jwtManager,
	)

	log.Println("✅ Routes configured")

	// Start server
	go startServer(app, cfg)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

func initRepositories(db *gorm.DB) *repository.Repositories {
	return &repository.Repositories{
		User:         mysql.NewUserRepository(db),
		Note:         mysql.NewNoteRepository(db),
		RefreshToken: mysql.NewRefreshTokenRepository(db),
	}
}

type Services struct {
	auth service.AuthService
	note service.NoteService
}

func initServices(repos *repository.Repositories, jwtManager *utils.JWTManager) *Services {
	return &Services{
		auth: impl.NewAuthService(repos.User, repos.RefreshToken, jwtManager),
		note: impl.NewNoteService(repos.Note, repos.User),
	}
}

type Handlers struct {
	auth   *handler.AuthHandler
	note   *handler.NoteHandler
	health *handler.HealthHandler
}

func initHandlers(services *Services, db *gorm.DB) *Handlers {
	return &Handlers{
		auth:   handler.NewAuthHandler(services.auth),
		note:   handler.NewNoteHandler(services.note),
		health: handler.NewHealthHandler(db),
	}
}

func setupMiddlewares(app *fiber.App, cfg *configs.Config) {
	// Request ID
	app.Use(middleware.RequestID())

	// Recover from panics
	app.Use(recover.New())

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path} | ${locals:requestID}\n",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// General rate limiting
	app.Use(middleware.GeneralRateLimit())
}

func startServer(app *fiber.App, cfg *configs.Config) {
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("✅ Server is running on http://localhost%s", addr)
	log.Printf("📚 API Documentation: http://localhost%s/api/v1/health", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
