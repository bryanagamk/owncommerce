package middleware

import (
	"strings"

	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
)

type TenantMiddleware struct {
	tenantSvc *tenant.Service
}

func NewTenantMiddleware(tenantSvc *tenant.Service) *TenantMiddleware {
	return &TenantMiddleware{tenantSvc: tenantSvc}
}

// ResolveTenant identifies tenant from X-Tenant-Slug (dev) or Host header.
func (m *TenantMiddleware) ResolveTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		if slug := strings.TrimSpace(c.Get("X-Tenant-Slug")); slug != "" {
			t, err := m.tenantSvc.ResolveBySlug(ctx, slug)
			if err != nil {
				return response.NotFound(c, "tenant not found")
			}
			c.Locals(LocalTenantKey, t)
			c.Locals(LocalTenantIDKey, &t.ID)
			return c.Next()
		}

		host := c.Hostname()
		if host == "" {
			return c.Next()
		}

		t, err := m.tenantSvc.ResolveByHost(ctx, host)
		if err != nil {
			// Allow platform routes without tenant context
			return c.Next()
		}

		c.Locals(LocalTenantKey, t)
		c.Locals(LocalTenantIDKey, &t.ID)
		return c.Next()
	}
}

// RequireTenant ensures a tenant is resolved before proceeding.
func (m *TenantMiddleware) RequireTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if TenantFromLocals(c) == nil {
			return response.BadRequest(c, "tenant context required (use X-Tenant-Slug header in local dev)")
		}
		return c.Next()
	}
}
