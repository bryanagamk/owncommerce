package category

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("category not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repository) Update(ctx context.Context, c *Category) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *Repository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&Category{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Category, error) {
	var c Category
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) FindBySlug(ctx context.Context, tenantID uuid.UUID, slug string) (*Category, error) {
	var c Category
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND slug = ?", tenantID, slug).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) List(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]Category, error) {
	var items []Category
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if activeOnly {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Order("sort_order ASC, name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) SlugExists(ctx context.Context, tenantID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&Category{}).Where("tenant_id = ? AND slug = ?", tenantID, slug)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
