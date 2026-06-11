package handler

import (
	"strconv"

	"github.com/agamtech/owncommerce/apps/api/internal/commerce/category"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/order"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/product"
	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MerchantHandler struct {
	tenantSvc   *tenant.Service
	categorySvc *category.Service
	productSvc  *product.Service
	orderSvc    *order.Service
}

func NewMerchantHandler(
	tenantSvc *tenant.Service,
	categorySvc *category.Service,
	productSvc *product.Service,
	orderSvc *order.Service,
) *MerchantHandler {
	return &MerchantHandler{
		tenantSvc:   tenantSvc,
		categorySvc: categorySvc,
		productSvc:  productSvc,
		orderSvc:    orderSvc,
	}
}

func (h *MerchantHandler) UpdateStoreSettings(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}

	var req struct {
		Name         *string `json:"name"`
		Description  *string `json:"description"`
		ContactEmail *string `json:"contact_email"`
		ContactPhone *string `json:"contact_phone"`
		Address      *string `json:"address"`
		City         *string `json:"city"`
		Province     *string `json:"province"`
		PostalCode   *string `json:"postal_code"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	t, err := h.tenantSvc.UpdateStore(c.UserContext(), *tenantID, tenant.UpdateStoreInput{
		Name: req.Name, Description: req.Description, ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone, Address: req.Address, City: req.City,
		Province: req.Province, PostalCode: req.PostalCode,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, t)
}

func (h *MerchantHandler) ListCategories(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	items, err := h.categorySvc.List(c.UserContext(), *tenantID, false)
	if err != nil {
		return response.InternalError(c, "failed to list categories")
	}
	return response.OK(c, items)
}

func (h *MerchantHandler) CreateCategory(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return response.BadRequest(c, "name is required")
	}
	item, err := h.categorySvc.Create(c.UserContext(), category.CreateInput{
		TenantID: *tenantID, Name: req.Name, Slug: req.Slug,
		Description: req.Description, IsActive: req.IsActive, SortOrder: req.SortOrder,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, item)
}

func (h *MerchantHandler) UpdateCategory(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid category id")
	}
	var req struct {
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		IsActive    *bool   `json:"is_active"`
		SortOrder   *int    `json:"sort_order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	item, err := h.categorySvc.Update(c.UserContext(), *tenantID, id, category.UpdateInput{
		Name: req.Name, Slug: req.Slug, Description: req.Description,
		IsActive: req.IsActive, SortOrder: req.SortOrder,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, item)
}

func (h *MerchantHandler) DeleteCategory(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid category id")
	}
	if err := h.categorySvc.Delete(c.UserContext(), *tenantID, id); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Message(c, "category deleted")
}

func (h *MerchantHandler) ListProducts(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
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
		TenantID: *tenantID, CategoryID: categoryID, Status: c.Query("status"),
		Search: c.Query("q"), Limit: queryInt(c, "limit", 20), Offset: queryInt(c, "offset", 0),
	})
	if err != nil {
		return response.InternalError(c, "failed to list products")
	}
	return response.Paginated(c, items, total, queryInt(c, "limit", 20), queryInt(c, "offset", 0))
}

func (h *MerchantHandler) CreateProduct(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	var req struct {
		CategoryID        *string `json:"category_id"`
		Name              string  `json:"name"`
		Slug              string  `json:"slug"`
		Description       string  `json:"description"`
		SKU               string  `json:"sku"`
		Price             int64   `json:"price"`
		Status            string  `json:"status"`
		IsFeatured        bool    `json:"is_featured"`
		Quantity          int     `json:"quantity"`
		LowStockThreshold int     `json:"low_stock_threshold"`
	}
	if err := c.BodyParser(&req); err != nil || req.Name == "" {
		return response.BadRequest(c, "name is required")
	}
	var categoryID *uuid.UUID
	if req.CategoryID != nil && *req.CategoryID != "" {
		id, parseErr := uuid.Parse(*req.CategoryID)
		if parseErr != nil {
			return response.BadRequest(c, "invalid category_id")
		}
		categoryID = &id
	}
	item, err := h.productSvc.Create(c.UserContext(), product.CreateInput{
		TenantID: *tenantID, CategoryID: categoryID, Name: req.Name, Slug: req.Slug,
		Description: req.Description, SKU: req.SKU, Price: req.Price, Status: req.Status,
		IsFeatured: req.IsFeatured, Quantity: req.Quantity, LowStockThreshold: req.LowStockThreshold,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, item)
}

func (h *MerchantHandler) GetProduct(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}
	item, err := h.productSvc.GetByID(c.UserContext(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "product not found")
	}
	return response.OK(c, item)
}

func (h *MerchantHandler) UpdateProduct(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}
	var req struct {
		CategoryID  *string `json:"category_id"`
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		SKU         *string `json:"sku"`
		Price       *int64  `json:"price"`
		Status      *string `json:"status"`
		IsFeatured  *bool   `json:"is_featured"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	var categoryID *uuid.UUID
	if req.CategoryID != nil && *req.CategoryID != "" {
		parsed, parseErr := uuid.Parse(*req.CategoryID)
		if parseErr != nil {
			return response.BadRequest(c, "invalid category_id")
		}
		categoryID = &parsed
	}
	item, err := h.productSvc.Update(c.UserContext(), *tenantID, id, product.UpdateInput{
		CategoryID: categoryID, Name: req.Name, Slug: req.Slug, Description: req.Description,
		SKU: req.SKU, Price: req.Price, Status: req.Status, IsFeatured: req.IsFeatured,
	})
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, item)
}

func (h *MerchantHandler) DeleteProduct(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}
	if err := h.productSvc.Delete(c.UserContext(), *tenantID, id); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Message(c, "product deleted")
}

func (h *MerchantHandler) UpdateInventory(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}
	var req struct {
		Quantity          *int `json:"quantity"`
		LowStockThreshold *int `json:"low_stock_threshold"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	inv, err := h.productSvc.UpdateInventory(c.UserContext(), *tenantID, id, req.Quantity, req.LowStockThreshold)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, inv)
}

func (h *MerchantHandler) UploadProductImage(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}
	file, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, "image file is required")
	}
	f, err := file.Open()
	if err != nil {
		return response.BadRequest(c, "failed to open image")
	}
	defer f.Close()

	sortOrder, _ := strconv.Atoi(c.FormValue("sort_order", "0"))
	img, err := h.productSvc.AddImage(c.UserContext(), *tenantID, id, file.Filename, f, sortOrder)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, img)
}

func (h *MerchantHandler) ListOrders(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	items, total, err := h.orderSvc.List(c.UserContext(), order.ListFilter{
		TenantID: *tenantID, Status: c.Query("status"),
		Limit: queryInt(c, "limit", 20), Offset: queryInt(c, "offset", 0),
	})
	if err != nil {
		return response.InternalError(c, "failed to list orders")
	}
	return response.Paginated(c, items, total, queryInt(c, "limit", 20), queryInt(c, "offset", 0))
}

func (h *MerchantHandler) GetOrder(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	item, err := h.orderSvc.GetByID(c.UserContext(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "order not found")
	}
	return response.OK(c, item)
}

func (h *MerchantHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	var req struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil || req.Status == "" {
		return response.BadRequest(c, "status is required")
	}
	item, err := h.orderSvc.UpdateStatus(c.UserContext(), *tenantID, id, req.Status, req.Note)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, item)
}

func (h *MerchantHandler) CancelOrder(c *fiber.Ctx) error {
	tenantID, err := merchantTenantID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	var req struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&req)
	item, err := h.orderSvc.Cancel(c.UserContext(), *tenantID, id, req.Note)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.OK(c, item)
}
