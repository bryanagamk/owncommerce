package middleware

import (
	"github.com/agamtech/owncommerce/apps/api/internal/core/auth"
	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	LocalUserKey    = "user"
	LocalClaimsKey  = "claims"
	LocalTenantKey  = "tenant"
	LocalTenantIDKey = "tenant_id"
)

type AuthUser struct {
	ID          uuid.UUID
	Email       string
	Name        string
	TenantID    *uuid.UUID
	Permissions []string
	Roles       []string
}

func UserFromLocals(c *fiber.Ctx) *AuthUser {
	if v := c.Locals(LocalUserKey); v != nil {
		if user, ok := v.(*AuthUser); ok {
			return user
		}
	}
	return nil
}

func ClaimsFromLocals(c *fiber.Ctx) *auth.Claims {
	if v := c.Locals(LocalClaimsKey); v != nil {
		if claims, ok := v.(*auth.Claims); ok {
			return claims
		}
	}
	return nil
}

func TenantFromLocals(c *fiber.Ctx) *tenant.Tenant {
	if v := c.Locals(LocalTenantKey); v != nil {
		if t, ok := v.(*tenant.Tenant); ok {
			return t
		}
	}
	return nil
}

func TenantIDFromLocals(c *fiber.Ctx) *uuid.UUID {
	if v := c.Locals(LocalTenantIDKey); v != nil {
		if id, ok := v.(*uuid.UUID); ok {
			return id
		}
	}
	return nil
}
