package cart

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("cart not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindActive(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, sessionID string) (*Cart, error) {
	var c Cart
	q := r.db.WithContext(ctx).Where("tenant_id = ? AND status = ?", tenantID, StatusActive)
	if customerID != nil {
		q = q.Where("customer_id = ?", *customerID)
	} else if sessionID != "" {
		q = q.Where("session_id = ?", sessionID)
	} else {
		return nil, ErrNotFound
	}
	if err := q.Preload("Items").First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Create(ctx context.Context, c *Cart) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repository) Save(ctx context.Context, c *Cart) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *Repository) FindItem(ctx context.Context, cartID, productID uuid.UUID) (*CartItem, error) {
	var item CartItem
	if err := r.db.WithContext(ctx).Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repository) FindItemByID(ctx context.Context, cartID, itemID uuid.UUID) (*CartItem, error) {
	var item CartItem
	if err := r.db.WithContext(ctx).Where("cart_id = ? AND id = ?", cartID, itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repository) UpsertItem(ctx context.Context, item *CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *Repository) DeleteItem(ctx context.Context, cartID, itemID uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("cart_id = ? AND id = ?", cartID, itemID).Delete(&CartItem{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ClearItems(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("cart_id = ?", cartID).Delete(&CartItem{}).Error
}

func (r *Repository) LoadWithProducts(ctx context.Context, cartID uuid.UUID) (*Cart, error) {
	var c Cart
	if err := r.db.WithContext(ctx).Preload("Items").First(&c, "id = ?", cartID).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
