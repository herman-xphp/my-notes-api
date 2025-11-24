package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

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
	// Load configuration
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	log.Printf("🚀 Starting %s in %s mode...", cfg.App.Name, cfg.App.Env)

	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer closeDatabase(db)

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExp,
		cfg.JWT.RefreshTokenExp,
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
		log.Printf("❌ Failed to close database: %v", err)
	}
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
