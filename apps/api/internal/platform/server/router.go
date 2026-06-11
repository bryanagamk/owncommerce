package server

import (
	"github.com/agamtech/owncommerce/apps/api/internal/handler"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Health     *handler.HealthHandler
	Auth       *handler.AuthHandler
	Tenant     *handler.TenantHandler
	Audit      *handler.AuditHandler
	IAM        *handler.IAMHandler
	Merchant   *handler.MerchantHandler
	Storefront *handler.StorefrontHandler
	Payment    *handler.PaymentHandler
	File       *handler.FileHandler
}

type Middlewares struct {
	Auth     *middleware.AuthMiddleware
	Tenant   *middleware.TenantMiddleware
	Customer *middleware.CustomerAuthMiddleware
}

func RegisterRoutes(app *fiber.App, h Handlers, m Middlewares) {
	app.Get("/health", h.Health.Health)
	app.Get("/files/*", h.File.Serve)

	v1 := app.Group("/v1")
	v1.Get("/health", h.Health.Health)
	v1.Post("/webhooks/midtrans", h.Payment.MidtransWebhook)

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

	merchant := v1.Group("/merchant", m.Auth.RequireAuth())
	merchant.Patch("/store/settings", m.Auth.RequirePermission("tenant.manage"), h.Merchant.UpdateStoreSettings)
	merchant.Get("/categories", m.Auth.RequirePermission("product.view"), h.Merchant.ListCategories)
	merchant.Post("/categories", m.Auth.RequirePermission("product.create"), h.Merchant.CreateCategory)
	merchant.Patch("/categories/:id", m.Auth.RequirePermission("product.update"), h.Merchant.UpdateCategory)
	merchant.Delete("/categories/:id", m.Auth.RequirePermission("product.delete"), h.Merchant.DeleteCategory)
	merchant.Get("/products", m.Auth.RequirePermission("product.view"), h.Merchant.ListProducts)
	merchant.Post("/products", m.Auth.RequirePermission("product.create"), h.Merchant.CreateProduct)
	merchant.Get("/products/:id", m.Auth.RequirePermission("product.view"), h.Merchant.GetProduct)
	merchant.Patch("/products/:id", m.Auth.RequirePermission("product.update"), h.Merchant.UpdateProduct)
	merchant.Delete("/products/:id", m.Auth.RequirePermission("product.delete"), h.Merchant.DeleteProduct)
	merchant.Patch("/products/:id/inventory", m.Auth.RequirePermission("product.update"), h.Merchant.UpdateInventory)
	merchant.Post("/products/:id/images", m.Auth.RequirePermission("product.update"), h.Merchant.UploadProductImage)
	merchant.Get("/orders", m.Auth.RequirePermission("order.view"), h.Merchant.ListOrders)
	merchant.Get("/orders/:id", m.Auth.RequirePermission("order.view"), h.Merchant.GetOrder)
	merchant.Patch("/orders/:id/status", m.Auth.RequirePermission("order.manage"), h.Merchant.UpdateOrderStatus)
	merchant.Post("/orders/:id/cancel", m.Auth.RequirePermission("order.manage"), h.Merchant.CancelOrder)

	storefront := v1.Group("/storefront", m.Tenant.ResolveTenant(), m.Tenant.RequireTenant(), m.Customer.Optional())
	storefront.Get("/home", h.Storefront.Home)
	storefront.Get("/categories", h.Storefront.ListCategories)
	storefront.Get("/products", h.Storefront.ListProducts)
	storefront.Get("/products/:slug", h.Storefront.GetProduct)
	storefront.Post("/auth/register", h.Storefront.RegisterCustomer)
	storefront.Post("/auth/login", h.Storefront.LoginCustomer)
	storefront.Get("/me", m.Customer.Require(), h.Storefront.CustomerMe)
	storefront.Get("/addresses", m.Customer.Require(), h.Storefront.ListAddresses)
	storefront.Post("/addresses", m.Customer.Require(), h.Storefront.CreateAddress)
	storefront.Get("/cart", h.Storefront.GetCart)
	storefront.Post("/cart/items", h.Storefront.AddCartItem)
	storefront.Patch("/cart/items/:id", h.Storefront.UpdateCartItem)
	storefront.Delete("/cart/items/:id", h.Storefront.RemoveCartItem)
	storefront.Post("/checkout", h.Storefront.Checkout)
	storefront.Post("/orders/:id/pay", h.Storefront.PayOrder)
	storefront.Get("/orders/:id", h.Storefront.GetOrder)
	storefront.Get("/orders", m.Customer.Require(), h.Storefront.ListOrders)

	tenantPublic := v1.Group("", m.Tenant.ResolveTenant())
	tenantPublic.Get("/store", h.Tenant.Resolve)
}
