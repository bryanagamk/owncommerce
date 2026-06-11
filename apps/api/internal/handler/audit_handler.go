package handler

import (
	"strconv"

	"github.com/agamtech/owncommerce/apps/api/internal/core/audit"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
	auditSvc *audit.Service
}

func NewAuditHandler(auditSvc *audit.Service) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc}
}

func (h *AuditHandler) List(c *fiber.Ctx) error {
	user := middleware.UserFromLocals(c)
	if user == nil {
		return response.Unauthorized(c, "authentication required")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	filter := audit.ListFilter{
		TenantID: user.TenantID,
		Action:   c.Query("action"),
		Limit:    limit,
		Offset:   offset,
	}

	logs, total, err := h.auditSvc.List(c.UserContext(), filter)
	if err != nil {
		return response.InternalError(c, "failed to load audit logs")
	}

	return response.Paginated(c, logs, total, limit, offset)
}
