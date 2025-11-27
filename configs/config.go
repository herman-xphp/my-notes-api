package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Security  SecurityConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type SecurityConfig struct {
	BcryptCons int
}

type RateLimitConfig struct {
	Enabled     bool
	MaxRequests int
	Window      time.Duration
}

type CORSConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "my-notes-api"),
			Port: getEnv("APP_PORT", "3000"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASS", ""),
			Name:            getEnv("DB_NAME", "my_notes_api"),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvAsInt("DB_OPEN_CONNS", 100),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 3600*time.Second),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", ""),
			AccessTokenDuration:  getEnvAsDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute),
			RefreshTokenDuration: getEnvAsDuration("JWT_REFRESH_TOKEN_DURATION", 7*24*time.Hour),
		},
		Security: SecurityConfig{
			BcryptCons: getEnvAsInt("BCRYPT_COST", 12),
		},
		RateLimit: RateLimitConfig{
			Enabled:     getEnvAsBool("RATE_LIMIT_ENABLED", true),
			MaxRequests: getEnvAsInt("RATE_LIMIT_MAX_REQUESTS", 100),
			Window:      getEnvAsDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CROS_ALLOWED_ORIGINS", "*"),
			AllowedMethods: getEnv("CROS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			AllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		},
	}

	// Validate configuration before return
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil

}

// Validate
func (c *Config) Validate() error {
	// Validate JWT Secret
	if err := c.validateJWTSecret(); err != nil {
		return err
	}

	// Validate database
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// VAlidate Bcrypt Cons (range 4-31, but < 10 too weak)
	if c.Security.BcryptCons < 10 || c.Security.BcryptCons > 31 {
		return fmt.Errorf("BCRYPT_COST must be between 10 and 31 (current: %d)", c.Security.BcryptCons)
	}

	return nil
}

// Validate JWTSecret
func (c *Config) validateJWTSecret() error {
	secret := c.JWT.Secret

	// Check 1: Secret not null
	if secret == "" {
		return fmt.Errorf("JWT_SECRET is required. Generate one using: go run scripts/generate-secret.go")
	}

	// Check 2: Secret not same to default/placeholder
	dangerousSecrets := []string{
		"CHANGE_THIS_TO_RANDOM_64_CHARS_STRING",
		"your-secret-key",
		"secret",
		"jwt-secret",
		"your_jwt_secret_here",
		"change-me",
		"changeme",
	}
	for _, dangerous := range dangerousSecrets {
		if secret == dangerous {
			return fmt.Errorf("JWT_SECRET is using a placeholder value. Generate a secure secret using: go run scripts/generate-secret.go")
		}
	}

	// Check 3: Minimum lenght (256 bits = 32 bytes = ~43 chars in base64)
	// Set minimum 32 chars for safety
	minLength := 32
	if len(secret) < minLength {
		return fmt.Errorf("JWT_SECRET is too short (minimum %d characters). Current lenght: %d. Generate a secure secret using: go run scripts/generate-secret.go", minLength, len(secret))
	}

	// Check 4: Warning for production with weak secret
	if c.App.Env == "production" && len(secret) < 64 {
		// This is just a warning, not an error
		fmt.Printf("WARNING: JWT_SECRET lenght is %d chars. For production, recommended minimum is 64 chars.\n", len(secret))
	}

	// Check 5: Secrets must not contain whitespace
	for _, char := range secret {
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			return fmt.Errorf("JWT_SECRET contains whitespace characters. Please remove them.")
		}
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
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

// Helper functions for parse environment variables
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {

		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
