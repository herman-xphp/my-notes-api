package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/herman-xphp/my-notes-api/pkg/response"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db        *gorm.DB
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Services  map[string]string `json:"services"`
}

// Check handles health check
func (h *HealthHandler) Check(c *fiber.Ctx) error {
	services := make(map[string]string)

	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		services["database"] = "unhealthy"
	} else if err := sqlDB.Ping(); err != nil {
		services["database"] = "unhealthy"
	} else {
		services["database"] = "healthy"
	}

	// Overall status
	status := "healthy"
	for _, serviceStatus := range services {
		if serviceStatus == "unhealthy" {
			status = "unhealthy"
			break
		}
	}

	// Build response
	healthResponse := HealthResponse{
		Status:    status,
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    time.Since(h.startTime).String(),
		Services:  services,
	}

	// Return appropriate status code
	if status == "unhealthy" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(healthResponse)
	}

	return response.Success(c, "Service is healthy", healthResponse)
}
