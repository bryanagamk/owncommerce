package middleware

import (
	"strings"

	"github.com/agamtech/owncommerce/apps/api/internal/core/auth"
	"github.com/agamtech/owncommerce/apps/api/internal/core/iam"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwt     *auth.JWTManager
	iamRepo *iam.Repository
}

func NewAuthMiddleware(jwt *auth.JWTManager, iamRepo *iam.Repository) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt, iamRepo: iamRepo}
}

func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractBearerToken(c.Get("Authorization"))
		if token == "" {
			return response.Unauthorized(c, "missing access token")
		}

		claims, err := m.jwt.ParseToken(token)
		if err != nil || claims.TokenType != auth.TokenTypeAccess {
			return response.Unauthorized(c, "invalid access token")
		}

		tenantID := claims.TenantID
		if headerTenant := c.Get("X-Tenant-ID"); headerTenant != "" && tenantID == nil {
			// optional override for platform users in dev
		}

		perms, roles, err := m.iamRepo.GetUserPermissions(c.UserContext(), claims.UserID, tenantID)
		if err != nil {
			return response.InternalError(c, "failed to load permissions")
		}

		c.Locals(LocalClaimsKey, claims)
		c.Locals(LocalUserKey, &AuthUser{
			ID:          claims.UserID,
			TenantID:    tenantID,
			Permissions: perms,
			Roles:       roles,
		})
		if tenantID != nil {
			c.Locals(LocalTenantIDKey, tenantID)
		}

		return c.Next()
	}
}

func (m *AuthMiddleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := UserFromLocals(c)
		if user == nil {
			return response.Unauthorized(c, "authentication required")
		}

		for _, p := range user.Permissions {
			if p == permission {
				return c.Next()
			}
		}

		return response.Forbidden(c, "insufficient permissions")
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
