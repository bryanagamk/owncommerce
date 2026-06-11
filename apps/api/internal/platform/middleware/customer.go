package middleware

import (
	"strings"

	"github.com/agamtech/owncommerce/apps/api/internal/core/auth"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const LocalCustomerIDKey = "customer_id"

type CustomerAuthMiddleware struct {
	jwt *auth.JWTManager
}

func NewCustomerAuthMiddleware(jwt *auth.JWTManager) *CustomerAuthMiddleware {
	return &CustomerAuthMiddleware{jwt: jwt}
}

func (m *CustomerAuthMiddleware) Optional() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractBearerToken(c.Get("Authorization"))
		if token == "" {
			return c.Next()
		}
		claims, err := m.jwt.ParseToken(token)
		if err != nil || claims.TokenType != auth.TokenTypeAccess || claims.ActorType != auth.ActorCustomer {
			return c.Next()
		}
		if claims.CustomerID != nil {
			c.Locals(LocalCustomerIDKey, *claims.CustomerID)
		}
		if claims.TenantID != nil {
			c.Locals(LocalTenantIDKey, claims.TenantID)
		}
		c.Locals(LocalClaimsKey, claims)
		return c.Next()
	}
}

func (m *CustomerAuthMiddleware) Require() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractBearerToken(c.Get("Authorization"))
		if token == "" {
			return response.Unauthorized(c, "missing access token")
		}
		claims, err := m.jwt.ParseToken(token)
		if err != nil || claims.TokenType != auth.TokenTypeAccess || claims.ActorType != auth.ActorCustomer {
			return response.Unauthorized(c, "invalid customer token")
		}
		if claims.CustomerID == nil || claims.TenantID == nil {
			return response.Unauthorized(c, "invalid customer token")
		}

		tenant := TenantFromLocals(c)
		if tenant != nil && tenant.ID != *claims.TenantID {
			return response.Forbidden(c, "tenant mismatch")
		}

		c.Locals(LocalCustomerIDKey, *claims.CustomerID)
		c.Locals(LocalTenantIDKey, claims.TenantID)
		c.Locals(LocalClaimsKey, claims)
		return c.Next()
	}
}

func CustomerIDFromLocals(c *fiber.Ctx) *uuid.UUID {
	if v := c.Locals(LocalCustomerIDKey); v != nil {
		if id, ok := v.(uuid.UUID); ok {
			return &id
		}
	}
	return nil
}

func CartSessionFromRequest(c *fiber.Ctx) string {
	return strings.TrimSpace(c.Get("X-Cart-Session"))
}
