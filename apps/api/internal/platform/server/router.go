package server

import (
	"github.com/agamtech/owncommerce/apps/api/internal/handler"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Health *handler.HealthHandler
	Auth   *handler.AuthHandler
	Tenant *handler.TenantHandler
	Audit  *handler.AuditHandler
	IAM    *handler.IAMHandler
}

type Middlewares struct {
	Auth   *middleware.AuthMiddleware
	Tenant *middleware.TenantMiddleware
}

func RegisterRoutes(app *fiber.App, h Handlers, m Middlewares) {
	app.Get("/health", h.Health.Health)

	v1 := app.Group("/v1")
	v1.Get("/health", h.Health.Health)

	authGroup := v1.Group("/auth")
	authGroup.Post("/register", h.Auth.Register)
	authGroup.Post("/login", h.Auth.Login)
	authGroup.Post("/refresh", h.Auth.Refresh)
	authGroup.Post("/logout", m.Auth.RequireAuth(), h.Auth.Logout)

	protected := v1.Group("", m.Auth.RequireAuth())
	protected.Get("/me", h.Auth.Me)
	protected.Get("/tenants/current", h.Tenant.Current)
	protected.Get("/audit-logs", m.Auth.RequirePermission("audit.view"), h.Audit.List)
	protected.Get("/iam/roles", m.Auth.RequirePermission("staff.manage"), h.IAM.ListRoles)
	protected.Get("/iam/permissions", m.Auth.RequirePermission("staff.manage"), h.IAM.ListPermissions)

	tenantPublic := v1.Group("", m.Tenant.ResolveTenant())
	tenantPublic.Get("/store", h.Tenant.Resolve)
}
