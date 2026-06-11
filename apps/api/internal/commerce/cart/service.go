package cart

import (
	"context"
	"errors"
	"fmt"

	"github.com/agamtech/owncommerce/apps/api/internal/commerce/product"
	"github.com/google/uuid"
)

type Service struct {
	repo        *Repository
	productRepo *product.Repository
}

func NewService(repo *Repository, productRepo *product.Repository) *Service {
	return &Service{repo: repo, productRepo: productRepo}
}

type CartView struct {
	Cart  *Cart
	Items []CartItemView
	Total int64
}

type CartItemView struct {
	CartItem
	ProductName string `json:"product_name"`
	ProductSlug string `json:"product_slug"`
	ImageURL    string `json:"image_url,omitempty"`
}

func (s *Service) GetOrCreate(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, sessionID string) (*Cart, error) {
	c, err := s.repo.FindActive(ctx, tenantID, customerID, sessionID)
	if err == nil {
		return c, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	c = &Cart{
		TenantID:   tenantID,
		CustomerID: customerID,
		SessionID:  sessionID,
		Status:     StatusActive,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetView(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, sessionID string) (*CartView, error) {
	c, err := s.GetOrCreate(ctx, tenantID, customerID, sessionID)
	if err != nil {
		return nil, err
	}
	return s.buildView(ctx, tenantID, c)
}

func (s *Service) AddItem(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, sessionID string, productID uuid.UUID, qty int) (*CartView, error) {
	if qty <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	p, err := s.productRepo.FindByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}
	if p.Status != product.StatusActive {
		return nil, fmt.Errorf("product is not available")
	}
	if p.Inventory != nil && p.Inventory.Quantity < qty {
		return nil, fmt.Errorf("insufficient stock")
	}

	c, err := s.GetOrCreate(ctx, tenantID, customerID, sessionID)
	if err != nil {
		return nil, err
	}

	item, err := s.repo.FindItem(ctx, c.ID, productID)
	if errors.Is(err, ErrNotFound) {
		item = &CartItem{
			CartID:    c.ID,
			ProductID: productID,
			Quantity:  qty,
			UnitPrice: p.Price,
		}
	} else if err != nil {
		return nil, err
	} else {
		item.Quantity += qty
		item.UnitPrice = p.Price
	}

	if err := s.repo.UpsertItem(ctx, item); err != nil {
		return nil, err
	}
	return s.buildView(ctx, tenantID, c)
}

func (s *Service) UpdateItem(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, sessionID string, itemID uuid.UUID, qty int) (*CartView, error) {
	c, err := s.GetOrCreate(ctx, tenantID, customerID, sessionID)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.FindItemByID(ctx, c.ID, itemID)
	if err != nil {
		return nil, err
	}
	if qty <= 0 {
		if err := s.repo.DeleteItem(ctx, c.ID, itemID); err != nil {
			return nil, err
		}
	} else {
		p, err := s.productRepo.FindByID(ctx, tenantID, item.ProductID)
		if err != nil {
			return nil, err
		}
		if p.Inventory != nil && p.Inventory.Quantity < qty {
			return nil, fmt.Errorf("insufficient stock")
		}
		item.Quantity = qty
		if err := s.repo.UpsertItem(ctx, item); err != nil {
			return nil, err
		}
	}
	return s.buildView(ctx, tenantID, c)
}

func (s *Service) RemoveItem(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, sessionID string, itemID uuid.UUID) (*CartView, error) {
	c, err := s.GetOrCreate(ctx, tenantID, customerID, sessionID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.DeleteItem(ctx, c.ID, itemID); err != nil {
		return nil, err
	}
	return s.buildView(ctx, tenantID, c)
}

func (s *Service) buildView(ctx context.Context, tenantID uuid.UUID, c *Cart) (*CartView, error) {
	c, err := s.repo.LoadWithProducts(ctx, c.ID)
	if err != nil {
		return nil, err
	}

	view := &CartView{Cart: c, Items: make([]CartItemView, 0, len(c.Items))}
	for _, item := range c.Items {
		p, err := s.productRepo.FindByID(ctx, tenantID, item.ProductID)
		if err != nil {
			continue
		}
		imgURL := ""
		if len(p.Images) > 0 {
			imgURL = p.Images[0].URL
		}
		view.Items = append(view.Items, CartItemView{
			CartItem:    item,
			ProductName: p.Name,
			ProductSlug: p.Slug,
			ImageURL:    imgURL,
		})
		view.Total += item.UnitPrice * int64(item.Quantity)
	}
	return view, nil
}

func (s *Service) Clear(ctx context.Context, cartID uuid.UUID) error {
	return s.repo.ClearItems(ctx, cartID)
}
