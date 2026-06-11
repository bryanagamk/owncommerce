package tenant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTenantNotFound = errors.New("tenant not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, tenant *Tenant, domain *TenantDomain) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}
		domain.TenantID = tenant.ID
		return tx.Create(domain).Error
	})
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) FindBySlug(ctx context.Context, slug string) (*Tenant, error) {
	var tenant Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "slug = ?", slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *Repository) FindByDomain(ctx context.Context, domain string) (*Tenant, error) {
	var tenantDomain TenantDomain
	if err := r.db.WithContext(ctx).Where("domain = ?", domain).First(&tenantDomain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return r.FindByID(ctx, tenantDomain.TenantID)
}

func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Tenant{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) DomainExists(ctx context.Context, domain string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&TenantDomain{}).Where("domain = ?", domain).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ListDomains(ctx context.Context, tenantID uuid.UUID) ([]TenantDomain, error) {
	var domains []TenantDomain
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}
