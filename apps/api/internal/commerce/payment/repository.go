package payment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("payment not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repository) Update(ctx context.Context, p *Payment) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *Repository) FindByOrderID(ctx context.Context, tenantID, orderID uuid.UUID) (*Payment, error) {
	var p Payment
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND order_id = ?", tenantID, orderID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) FindByMidtransOrderID(ctx context.Context, midtransOrderID string) (*Payment, error) {
	var p Payment
	if err := r.db.WithContext(ctx).Where("midtrans_order_id = ?", midtransOrderID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}
