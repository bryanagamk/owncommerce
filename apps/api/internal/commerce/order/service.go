package order

import (
	"context"
	"fmt"
	"time"

	"github.com/agamtech/owncommerce/apps/api/internal/commerce/cart"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/product"
	"github.com/google/uuid"
)

type Service struct {
	repo        *Repository
	cartRepo    *cart.Repository
	productRepo *product.Repository
}

func NewService(repo *Repository, cartRepo *cart.Repository, productRepo *product.Repository) *Service {
	return &Service{repo: repo, cartRepo: cartRepo, productRepo: productRepo}
}

type CheckoutInput struct {
	TenantID        uuid.UUID
	CustomerID      *uuid.UUID
	SessionID       string
	RecipientName   string
	RecipientPhone  string
	ShippingAddress string
	ShippingCity    string
	ShippingProvince string
	ShippingPostal  string
	CustomerEmail   string
	CustomerNote    string
	ShippingCost    int64
}

func (s *Service) Checkout(ctx context.Context, input CheckoutInput) (*Order, error) {
	c, err := s.cartRepo.FindActive(ctx, input.TenantID, input.CustomerID, input.SessionID)
	if err != nil {
		return nil, fmt.Errorf("cart is empty")
	}
	c, err = s.cartRepo.LoadWithProducts(ctx, c.ID)
	if err != nil || len(c.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	var subtotal int64
	orderItems := make([]OrderItem, 0, len(c.Items))
	for _, item := range c.Items {
		p, err := s.productRepo.FindByID(ctx, input.TenantID, item.ProductID)
		if err != nil {
			return nil, err
		}
		if p.Status != product.StatusActive {
			return nil, fmt.Errorf("product %s is not available", p.Name)
		}
		if p.Inventory != nil && p.Inventory.Quantity < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s", p.Name)
		}
		lineSubtotal := p.Price * int64(item.Quantity)
		subtotal += lineSubtotal
		orderItems = append(orderItems, OrderItem{
			ProductID:   p.ID,
			ProductName: p.Name,
			SKU:         p.SKU,
			Quantity:    item.Quantity,
			UnitPrice:   p.Price,
			Subtotal:    lineSubtotal,
		})
	}

	total := subtotal + input.ShippingCost
	orderNumber := generateOrderNumber()

	o := &Order{
		TenantID:         input.TenantID,
		CustomerID:       input.CustomerID,
		OrderNumber:      orderNumber,
		Status:           StatusPendingPayment,
		PaymentStatus:    "pending",
		Subtotal:         subtotal,
		ShippingCost:     input.ShippingCost,
		Total:            total,
		RecipientName:    input.RecipientName,
		RecipientPhone:   input.RecipientPhone,
		ShippingAddress:  input.ShippingAddress,
		ShippingCity:     input.ShippingCity,
		ShippingProvince: input.ShippingProvince,
		ShippingPostal:   input.ShippingPostal,
		CustomerEmail:    input.CustomerEmail,
		CustomerNote:     input.CustomerNote,
	}

	if err := s.repo.Create(ctx, o, orderItems); err != nil {
		return nil, err
	}

	_ = s.cartRepo.ClearItems(ctx, c.ID)
	return s.repo.FindByID(ctx, input.TenantID, o.ID)
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Order, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]Order, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status, note string) (*Order, error) {
	return s.repo.UpdateStatus(ctx, tenantID, id, status, note)
}

func (s *Service) Cancel(ctx context.Context, tenantID, id uuid.UUID, note string) (*Order, error) {
	o, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if o.Status == StatusCompleted || o.Status == StatusCancelled {
		return nil, fmt.Errorf("order cannot be cancelled")
	}
	return s.repo.UpdateStatus(ctx, tenantID, id, StatusCancelled, note)
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%s-%s", time.Now().UTC().Format("20060102"), uuid.New().String()[:8])
}
