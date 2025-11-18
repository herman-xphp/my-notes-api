package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DBConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret          string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

type CORSConfig struct {
	AllowedOrigins string
}

func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// Parse JWT expiration
	accessExp, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_EXP", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_EXP: %w", err)
	}

	refreshExp, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_EXP", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TOKEN_EXP: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "my-notes-api"),
			Port: getEnv("APP_PORT", "3000"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASS", ""),
			Name:     getEnv("DB_NAME", "my_notes_api"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			AccessTokenExp:  accessExp,
			RefreshTokenExp: refreshExp,
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CROS_ALLOWED_ORIGINS", "*"),
		},
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	return nil
}

func (c *Config) GetDBDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
