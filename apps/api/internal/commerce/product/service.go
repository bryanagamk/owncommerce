package product

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/agamtech/owncommerce/apps/api/internal/platform/slug"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/storage"
	"github.com/google/uuid"
)

type Service struct {
	repo    *Repository
	storage *storage.LocalStorage
	baseURL string
}

func NewService(repo *Repository, store *storage.LocalStorage, baseURL string) *Service {
	return &Service{repo: repo, storage: store, baseURL: baseURL}
}

type CreateInput struct {
	TenantID          uuid.UUID
	CategoryID        *uuid.UUID
	Name              string
	Slug              string
	Description       string
	SKU               string
	Price             int64
	Status            string
	IsFeatured        bool
	Quantity          int
	LowStockThreshold int
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*Product, error) {
	prodSlug := input.Slug
	if prodSlug == "" {
		prodSlug = slug.FromName(input.Name)
	}
	exists, err := s.repo.SlugExists(ctx, input.TenantID, prodSlug, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("product slug already exists")
	}

	status := input.Status
	if status == "" {
		status = StatusDraft
	}

	p := &Product{
		TenantID:    input.TenantID,
		CategoryID:  input.CategoryID,
		Name:        input.Name,
		Slug:        prodSlug,
		Description: input.Description,
		SKU:         input.SKU,
		Price:       input.Price,
		Status:      status,
		IsFeatured:  input.IsFeatured,
	}
	inv := &Inventory{
		Quantity:          input.Quantity,
		LowStockThreshold: input.LowStockThreshold,
	}
	if inv.LowStockThreshold == 0 {
		inv.LowStockThreshold = 5
	}

	if err := s.repo.Create(ctx, p, inv); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, input.TenantID, p.ID)
}

type UpdateInput struct {
	CategoryID  *uuid.UUID
	Name        *string
	Slug        *string
	Description *string
	SKU         *string
	Price       *int64
	Status      *string
	IsFeatured  *bool
}

func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, input UpdateInput) (*Product, error) {
	p, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if input.CategoryID != nil {
		p.CategoryID = input.CategoryID
	}
	if input.Name != nil {
		p.Name = *input.Name
	}
	if input.Slug != nil {
		exists, err := s.repo.SlugExists(ctx, tenantID, *input.Slug, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("product slug already exists")
		}
		p.Slug = *input.Slug
	}
	if input.Description != nil {
		p.Description = *input.Description
	}
	if input.SKU != nil {
		p.SKU = *input.SKU
	}
	if input.Price != nil {
		p.Price = *input.Price
	}
	if input.Status != nil {
		p.Status = *input.Status
	}
	if input.IsFeatured != nil {
		p.IsFeatured = *input.IsFeatured
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Product, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s *Service) GetBySlug(ctx context.Context, tenantID uuid.UUID, prodSlug string, activeOnly bool) (*Product, error) {
	return s.repo.FindBySlug(ctx, tenantID, prodSlug, activeOnly)
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]Product, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *Service) UpdateInventory(ctx context.Context, tenantID, productID uuid.UUID, quantity, lowStock *int) (*Inventory, error) {
	return s.repo.UpdateInventory(ctx, tenantID, productID, quantity, lowStock)
}

func (s *Service) AddImage(ctx context.Context, tenantID, productID uuid.UUID, filename string, reader io.Reader, sortOrder int) (*ProductImage, error) {
	if _, err := s.repo.FindByID(ctx, tenantID, productID); err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".jpg"
	}
	safeName := uuid.New().String() + ext
	relPath, err := s.storage.Save(tenantID, "products", safeName, reader)
	if err != nil {
		return nil, err
	}

	img := &ProductImage{
		ProductID: productID,
		Path:      relPath,
		URL:       s.fileURL(relPath),
		SortOrder: sortOrder,
	}
	if err := s.repo.AddImage(ctx, img); err != nil {
		return nil, err
	}
	return img, nil
}

func (s *Service) fileURL(relPath string) string {
	if s.baseURL == "" {
		return "/files/" + relPath
	}
	return strings.TrimRight(s.baseURL, "/") + "/files/" + relPath
}