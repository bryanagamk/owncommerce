package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("product not found")

type ListFilter struct {
	TenantID   uuid.UUID
	CategoryID *uuid.UUID
	Status     string
	Search     string
	Featured   *bool
	Limit      int
	Offset     int
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *Product, inv *Inventory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		inv.ProductID = p.ID
		inv.TenantID = p.TenantID
		return tx.Create(inv).Error
	})
}

func (r *Repository) Update(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&Product{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Product, error) {
	var p Product
	if err := r.db.WithContext(ctx).
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("Inventory").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) FindBySlug(ctx context.Context, tenantID uuid.UUID, prodSlug string, activeOnly bool) (*Product, error) {
	var p Product
	q := r.db.WithContext(ctx).
		Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("Inventory").
		Where("tenant_id = ? AND slug = ?", tenantID, prodSlug)
	if activeOnly {
		q = q.Where("status = ?", StatusActive)
	}
	if err := q.First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) List(ctx context.Context, filter ListFilter) ([]Product, int64, error) {
	q := r.db.WithContext(ctx).Model(&Product{}).Where("tenant_id = ?", filter.TenantID)
	if filter.CategoryID != nil {
		q = q.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Featured != nil {
		q = q.Where("is_featured = ?", *filter.Featured)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", like, like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var items []Product
	if err := q.Preload("Images", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("Inventory").
		Order("created_at DESC").
		Limit(limit).Offset(filter.Offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) SlugExists(ctx context.Context, tenantID uuid.UUID, prodSlug string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&Product{}).Where("tenant_id = ? AND slug = ?", tenantID, prodSlug)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) AddImage(ctx context.Context, img *ProductImage) error {
	return r.db.WithContext(ctx).Create(img).Error
}

func (r *Repository) UpdateInventory(ctx context.Context, tenantID, productID uuid.UUID, quantity, lowStock *int) (*Inventory, error) {
	var inv Inventory
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND product_id = ?", tenantID, productID).First(&inv).Error; err != nil {
		return nil, err
	}
	if quantity != nil {
		inv.Quantity = *quantity
	}
	if lowStock != nil {
		inv.LowStockThreshold = *lowStock
	}
	if err := r.db.WithContext(ctx).Save(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *Repository) DecrementStock(ctx context.Context, tenantID, productID uuid.UUID, qty int) error {
	res := r.db.WithContext(ctx).Model(&Inventory{}).
		Where("tenant_id = ? AND product_id = ? AND quantity >= ?", tenantID, productID, qty).
		Update("quantity", gorm.Expr("quantity - ?", qty))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("insufficient stock")
	}
	return nil
}
