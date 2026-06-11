package handler

import (
	"github.com/agamtech/owncommerce/apps/api/internal/platform/database"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	status := "ok"
	dbStatus := "ok"

	if err := database.Ping(h.db); err != nil {
		status = "degraded"
		dbStatus = "down"
	}

	return response.OK(c, fiber.Map{
		"status":   status,
		"database": dbStatus,
		"service":  "owncommerce-api",
	})
}
