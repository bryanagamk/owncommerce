package customer

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("customer not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *Customer) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*Customer, error) {
	var c Customer
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND email = ?", tenantID, email).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Customer, error) {
	var c Customer
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Update(ctx context.Context, c *Customer) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *Repository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&Customer{}).Where("id = ?", id).Update("last_login_at", now).Error
}

func (r *Repository) CreateAddress(ctx context.Context, a *Address) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if a.IsDefault {
			if err := tx.Model(&Address{}).
				Where("tenant_id = ? AND customer_id = ?", a.TenantID, a.CustomerID).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(a).Error
	})
}

func (r *Repository) UpdateAddress(ctx context.Context, a *Address) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if a.IsDefault {
			if err := tx.Model(&Address{}).
				Where("tenant_id = ? AND customer_id = ? AND id != ?", a.TenantID, a.CustomerID, a.ID).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(a).Error
	})
}

func (r *Repository) DeleteAddress(ctx context.Context, tenantID, customerID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("tenant_id = ? AND customer_id = ? AND id = ?", tenantID, customerID, id).Delete(&Address{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) FindAddress(ctx context.Context, tenantID, customerID, id uuid.UUID) (*Address, error) {
	var a Address
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND customer_id = ? AND id = ?", tenantID, customerID, id).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) ListAddresses(ctx context.Context, tenantID, customerID uuid.UUID) ([]Address, error) {
	var items []Address
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Order("is_default DESC, created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
