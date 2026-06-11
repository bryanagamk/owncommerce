package handler

import (
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func merchantTenantID(c *fiber.Ctx) (*uuid.UUID, error) {
	user := middleware.UserFromLocals(c)
	if user == nil || user.TenantID == nil {
		return nil, response.Forbidden(c, "merchant tenant required")
	}
	return user.TenantID, nil
}

func storefrontTenantID(c *fiber.Ctx) (uuid.UUID, error) {
	t := middleware.TenantFromLocals(c)
	if t == nil {
		return uuid.Nil, response.BadRequest(c, "tenant context required (use X-Tenant-Slug header in local dev)")
	}
	return t.ID, nil
}

func queryInt(c *fiber.Ctx, key string, fallback int) int {
	v := c.QueryInt(key, fallback)
	if v < 0 {
		return fallback
	}
	return v
}
