package handler

import (
	"github.com/agamtech/owncommerce/apps/api/internal/core/iam"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
)

type IAMHandler struct {
	iamRepo *iam.Repository
}

func NewIAMHandler(iamRepo *iam.Repository) *IAMHandler {
	return &IAMHandler{iamRepo: iamRepo}
}

func (h *IAMHandler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.iamRepo.ListRoles(c.UserContext())
	if err != nil {
		return response.InternalError(c, "failed to load roles")
	}
	return response.OK(c, roles)
}

func (h *IAMHandler) ListPermissions(c *fiber.Ctx) error {
	perms, err := h.iamRepo.ListPermissions(c.UserContext())
	if err != nil {
		return response.InternalError(c, "failed to load permissions")
	}
	return response.OK(c, perms)
}
