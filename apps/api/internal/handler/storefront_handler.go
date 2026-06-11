package handler

import (
	"errors"

	"github.com/agamtech/owncommerce/apps/api/internal/commerce/cart"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/category"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/customer"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/order"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/payment"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/product"
	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/middleware"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StorefrontHandler struct {
	tenantSvc   *tenant.Service
	categorySvc *category.Service
	productSvc  *product.Service
	customerSvc *customer.Service
	cartSvc     *cart.Service
	orderSvc    *order.Service
	paymentSvc  *payment.Service
}

func NewStorefrontHandler(
	tenantSvc *tenant.Service,
	categorySvc *category.Service,
	productSvc *product.Service,
	customerSvc *customer.Service,
	cartSvc *cart.Service,
	orderSvc *order.Service,
	paymentSvc *payment.Service,
) *StorefrontHandler {
	return &StorefrontHandler{
		tenantSvc: tenantSvc, categorySvc: categorySvc, productSvc: productSvc,
		customerSvc: customerSvc, cartSvc: cartSvc, orderSvc: orderSvc, paymentSvc: paymentSvc,
	}
}

func (h *StorefrontHandler) Home(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	t, err := h.tenantSvc.GetByID(c.UserContext(), tenantID)
	if err != nil {
		return response.NotFound(c, "store not found")
	}
	featured := true
	products, _, err := h.productSvc.List(c.UserContext(), product.ListFilter{
		TenantID: tenantID, Status: product.StatusActive, Featured: &featured, Limit: 8,
	})
	if err != nil {
		return response.InternalError(c, "failed to load home")
	}
	return response.OK(c, fiber.Map{"store": t, "featured_products": products})
}

func (h *StorefrontHandler) ListCategories(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	items, err := h.categorySvc.List(c.UserContext(), tenantID, true)
	if err != nil {
		return response.InternalError(c, "failed to list categories")
	}
	return response.OK(c, items)
}

func (h *StorefrontHandler) ListProducts(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	var categoryID *uuid.UUID
	if v := c.Query("category_id"); v != "" {
		id, parseErr := uuid.Parse(v)
		if parseErr != nil {
			return response.BadRequest(c, "invalid category_id")
		}
		categoryID = &id
	}
	items, total, err := h.productSvc.List(c.UserContext(), product.ListFilter{
		TenantID: tenantID, CategoryID: categoryID, Status: product.StatusActive,
		Search: c.Query("q"), Limit: queryInt(c, "limit", 20), Offset: queryInt(c, "offset", 0),
	})
	if err != nil {
		return response.InternalError(c, "failed to list products")
	}
	return response.Paginated(c, items, total, queryInt(c, "limit", 20), queryInt(c, "offset", 0))
}

func (h *StorefrontHandler) GetProduct(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	item, err := h.productSvc.GetBySlug(c.UserContext(), tenantID, c.Params("slug"), true)
	if err != nil {
		return response.NotFound(c, "product not found")
	}
	return response.OK(c, item)
}

func (h *StorefrontHandler) RegisterCustomer(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
	}
	if err := c.BodyParser(&req); err != nil || req.Email == "" || req.Password == "" || req.Name == "" {
		return response.BadRequest(c, "email, password, and name are required")
	}
	result, err := h.customerSvc.Register(c.UserContext(), customer.RegisterInput{
		TenantID: tenantID, Email: req.Email, Password: req.Password, Name: req.Name, Phone: req.Phone,
	})
	if err != nil {
		if errors.Is(err, customer.ErrEmailTaken) {
			return response.BadRequest(c, err.Error())
		}
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, fiber.Map{
		"customer": fiber.Map{"id": result.Customer.ID, "email": result.Customer.Email, "name": result.Customer.Name},
		"tokens":   result.Tokens,
	})
}

func (h *StorefrontHandler) LoginCustomer(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	result, err := h.customerSvc.Login(c.UserContext(), customer.LoginInput{
		TenantID: tenantID, Email: req.Email, Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, customer.ErrInvalidCredentials) {
			return response.Unauthorized(c, err.Error())
		}
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, fiber.Map{
		"customer": fiber.Map{"id": result.Customer.ID, "email": result.Customer.Email, "name": result.Customer.Name},
		"tokens":   result.Tokens,
	})
}

func (h *StorefrontHandler) CustomerMe(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	customerID := middleware.CustomerIDFromLocals(c)
	if customerID == nil {
		return response.Unauthorized(c, "authentication required")
	}
	cust, err := h.customerSvc.GetByID(c.UserContext(), tenantID, *customerID)
	if err != nil {
		return response.NotFound(c, "customer not found")
	}
	return response.OK(c, cust)
}

func (h *StorefrontHandler) UpdateCustomerMe(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	customerID := middleware.CustomerIDFromLocals(c)
	if customerID == nil {
		return response.Unauthorized(c, "authentication required")
	}
	var req struct {
		Name  *string `json:"name"`
		Phone *string `json:"phone"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == nil && req.Phone == nil {
		return response.BadRequest(c, "name or phone is required")
	}
	cust, err := h.customerSvc.UpdateProfile(c.UserContext(), tenantID, *customerID, customer.UpdateProfileInput{
		Name: req.Name, Phone: req.Phone,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, cust)
}

func (h *StorefrontHandler) ListAddresses(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	customerID := middleware.CustomerIDFromLocals(c)
	if customerID == nil {
		return response.Unauthorized(c, "authentication required")
	}
	items, err := h.customerSvc.ListAddresses(c.UserContext(), tenantID, *customerID)
	if err != nil {
		return response.InternalError(c, "failed to list addresses")
	}
	return response.OK(c, items)
}

func (h *StorefrontHandler) CreateAddress(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	customerID := middleware.CustomerIDFromLocals(c)
	if customerID == nil {
		return response.Unauthorized(c, "authentication required")
	}
	var req customer.AddressInput
	if err := c.BodyParser(&req); err != nil || req.RecipientName == "" {
		return response.BadRequest(c, "invalid address data")
	}
	item, err := h.customerSvc.CreateAddress(c.UserContext(), tenantID, *customerID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, item)
}

func (h *StorefrontHandler) UpdateAddress(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	customerID := middleware.CustomerIDFromLocals(c)
	if customerID == nil {
		return response.Unauthorized(c, "authentication required")
	}
	addressID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid address id")
	}
	var req customer.AddressInput
	if err := c.BodyParser(&req); err != nil || req.RecipientName == "" {
		return response.BadRequest(c, "invalid address data")
	}
	item, err := h.customerSvc.UpdateAddress(c.UserContext(), tenantID, *customerID, addressID, req)
	if err != nil {
		return response.NotFound(c, "address not found")
	}
	return response.OK(c, item)
}

func (h *StorefrontHandler) DeleteAddress(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	customerID := middleware.CustomerIDFromLocals(c)
	if customerID == nil {
		return response.Unauthorized(c, "authentication required")
	}
	addressID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid address id")
	}
	if err := h.customerSvc.DeleteAddress(c.UserContext(), tenantID, *customerID, addressID); err != nil {
		return response.NotFound(c, "address not found")
	}
	return response.OK(c, fiber.Map{"deleted": true})
}

func (h *StorefrontHandler) GetCart(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	view, err := h.cartSvc.GetView(c.UserContext(), tenantID, middleware.CustomerIDFromLocals(c), middleware.CartSessionFromRequest(c))
	if err != nil {
		return response.InternalError(c, "failed to load cart")
	}
	return response.OK(c, view)
}

func (h *StorefrontHandler) AddCartItem(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return response.BadRequest(c, "invalid product_id")
	}
	qty := req.Quantity
	if qty <= 0 {
		qty = 1
	}
	view, err := h.cartSvc.AddItem(c.UserContext(), tenantID, middleware.CustomerIDFromLocals(c), middleware.CartSessionFromRequest(c), productID, qty)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, view)
}

func (h *StorefrontHandler) UpdateCartItem(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid item id")
	}
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	view, err := h.cartSvc.UpdateItem(c.UserContext(), tenantID, middleware.CustomerIDFromLocals(c), middleware.CartSessionFromRequest(c), itemID, req.Quantity)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, view)
}

func (h *StorefrontHandler) RemoveCartItem(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid item id")
	}
	view, err := h.cartSvc.RemoveItem(c.UserContext(), tenantID, middleware.CustomerIDFromLocals(c), middleware.CartSessionFromRequest(c), itemID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, view)
}

func (h *StorefrontHandler) Checkout(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	var req struct {
		RecipientName    string `json:"recipient_name"`
		RecipientPhone   string `json:"recipient_phone"`
		ShippingAddress  string `json:"shipping_address"`
		ShippingCity     string `json:"shipping_city"`
		ShippingProvince string `json:"shipping_province"`
		ShippingPostal   string `json:"shipping_postal_code"`
		CustomerEmail    string `json:"customer_email"`
		CustomerNote     string `json:"customer_note"`
		ShippingCost     int64  `json:"shipping_cost"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.RecipientName == "" || req.RecipientPhone == "" || req.ShippingAddress == "" {
		return response.BadRequest(c, "shipping information is required")
	}

	o, err := h.orderSvc.Checkout(c.UserContext(), order.CheckoutInput{
		TenantID: tenantID, CustomerID: middleware.CustomerIDFromLocals(c),
		SessionID: middleware.CartSessionFromRequest(c),
		RecipientName: req.RecipientName, RecipientPhone: req.RecipientPhone,
		ShippingAddress: req.ShippingAddress, ShippingCity: req.ShippingCity,
		ShippingProvince: req.ShippingProvince, ShippingPostal: req.ShippingPostal,
		CustomerEmail: req.CustomerEmail, CustomerNote: req.CustomerNote, ShippingCost: req.ShippingCost,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, o)
}

func (h *StorefrontHandler) PayOrder(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	result, err := h.paymentSvc.CreateSnapPayment(c.UserContext(), tenantID, orderID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, result)
}

func (h *StorefrontHandler) GetOrder(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	orderID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	o, err := h.orderSvc.GetByID(c.UserContext(), tenantID, orderID)
	if err != nil {
		return response.NotFound(c, "order not found")
	}
	if cid := middleware.CustomerIDFromLocals(c); cid != nil && o.CustomerID != nil && *o.CustomerID != *cid {
		return response.Forbidden(c, "access denied")
	}
	return response.OK(c, o)
}

func (h *StorefrontHandler) ListOrders(c *fiber.Ctx) error {
	tenantID, err := storefrontTenantID(c)
	if err != nil {
		return err
	}
	customerID := middleware.CustomerIDFromLocals(c)
	if customerID == nil {
		return response.Unauthorized(c, "authentication required")
	}
	items, total, err := h.orderSvc.List(c.UserContext(), order.ListFilter{
		TenantID: tenantID, CustomerID: customerID,
		Limit: queryInt(c, "limit", 20), Offset: queryInt(c, "offset", 0),
	})
	if err != nil {
		return response.InternalError(c, "failed to list orders")
	}
	return response.Paginated(c, items, total, queryInt(c, "limit", 20), queryInt(c, "offset", 0))
}
