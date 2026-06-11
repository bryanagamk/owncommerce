package handler

import (
	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
)

type TenantHandler struct {
	tenantSvc *tenant.Service
	tenantRepo *tenant.Repository
}

func NewTenantHandler(tenantSvc *tenant.Service, tenantRepo *tenant.Repository) *TenantHandler {
	return &TenantHandler{tenantSvc: tenantSvc, tenantRepo: tenantRepo}
}

func (h *TenantHandler) Current(c *fiber.Ctx) error {
	user := middleware.UserFromLocals(c)
	if user == nil || user.TenantID == nil {
		return response.NotFound(c, "tenant not found for user")
	}

	t, err := h.tenantSvc.GetByID(c.UserContext(), *user.TenantID)
	if err != nil {
		return response.NotFound(c, "tenant not found")
	}

	domains, err := h.tenantRepo.ListDomains(c.UserContext(), t.ID)
	if err != nil {
		return response.InternalError(c, "failed to load domains")
	}

	return response.OK(c, fiber.Map{
		"tenant":  t,
		"domains": domains,
	})
}

func (h *TenantHandler) Resolve(c *fiber.Ctx) error {
	t := middleware.TenantFromLocals(c)
	if t == nil {
		return response.NotFound(c, "tenant not found")
	}

	domains, err := h.tenantRepo.ListDomains(c.UserContext(), t.ID)
	if err != nil {
		return response.InternalError(c, "failed to load domains")
	}

	return response.OK(c, fiber.Map{
		"tenant":  t,
		"domains": domains,
	})
}
