package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/herman-xphp/my-notes-api/configs"
)

func Connect(cfg *configs.DatabaseConfig) (*gorm.DB, error) {
	// Note: We don't have IsProduction check here easily without AppConfig,
	// but we can default to Info or pass it in.
	// For now, let's just use Info as per original NewMySQL default or
	// we could infer from environment if needed, but keeping it simple as per request.
	// Actually, the original code had logic for log level based on env.
	// Let's keep it simple for now and use NewMySQL.

	return NewMySQL(Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		),
		MaxIdleConns: cfg.MaxIdleConns,
		MaxOpenConns: cfg.MaxOpenConns,
		MaxLifeTime:  cfg.ConnMaxLifetime,
		LogLevel:     logger.Info, // Default to Info
	})
}

type Config struct {
	DSN          string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifeTime  time.Duration
	LogLevel     logger.LogLevel
}

func NewMySQL(cfg Config) (*gorm.DB, error) {
	// Set defaults
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 100
	}
	if cfg.MaxLifeTime == 0 {
		cfg.MaxLifeTime = time.Hour
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Connect to database
	db, err := gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection successfully")

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}
