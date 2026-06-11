package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("order not found")

type ListFilter struct {
	TenantID   uuid.UUID
	CustomerID *uuid.UUID
	Status     string
	Limit      int
	Offset     int
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, o *Order, items []OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = o.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		h := OrderStatusHistory{OrderID: o.ID, Status: o.Status, Note: "order created"}
		return tx.Create(&h).Error
	})
}

func (r *Repository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Order, error) {
	var o Order
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("History", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&o).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repository) FindByOrderNumber(ctx context.Context, tenantID uuid.UUID, orderNumber string) (*Order, error) {
	var o Order
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Where("tenant_id = ? AND order_number = ?", tenantID, orderNumber).
		First(&o).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repository) List(ctx context.Context, filter ListFilter) ([]Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&Order{}).Where("tenant_id = ?", filter.TenantID)
	if filter.CustomerID != nil {
		q = q.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var items []Order
	if err := q.Preload("Items").
		Order("created_at DESC").
		Limit(limit).Offset(filter.Offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status, note string) (*Order, error) {
	var o Order
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("tenant_id = ? AND id = ?", tenantID, id).First(&o).Error; err != nil {
			return err
		}
		o.Status = status
		if err := tx.Save(&o).Error; err != nil {
			return err
		}
		return tx.Create(&OrderStatusHistory{OrderID: o.ID, Status: status, Note: note}).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return r.FindByID(ctx, tenantID, id)
}

func (r *Repository) MarkPaid(ctx context.Context, tenantID, id uuid.UUID) error {
	now := gorm.Expr("NOW()")
	return r.db.WithContext(ctx).Model(&Order{}).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Updates(map[string]interface{}{
			"status":         StatusPaid,
			"payment_status": "paid",
			"paid_at":        now,
		}).Error
}
