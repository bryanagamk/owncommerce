package category

import (
	"context"
	"fmt"

	"github.com/agamtech/owncommerce/apps/api/internal/platform/slug"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	TenantID    uuid.UUID
	Name        string
	Slug        string
	Description string
	IsActive    bool
	SortOrder   int
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*Category, error) {
	catSlug := input.Slug
	if catSlug == "" {
		catSlug = slug.FromName(input.Name)
	}
	exists, err := s.repo.SlugExists(ctx, input.TenantID, catSlug, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("category slug already exists")
	}

	c := &Category{
		TenantID:    input.TenantID,
		Name:        input.Name,
		Slug:        catSlug,
		Description: input.Description,
		IsActive:    input.IsActive,
		SortOrder:   input.SortOrder,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

type UpdateInput struct {
	Name        *string
	Slug        *string
	Description *string
	IsActive    *bool
	SortOrder   *int
}

func (s *Service) Update(ctx context.Context, tenantID, id uuid.UUID, input UpdateInput) (*Category, error) {
	c, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		c.Name = *input.Name
	}
	if input.Slug != nil {
		exists, err := s.repo.SlugExists(ctx, tenantID, *input.Slug, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("category slug already exists")
		}
		c.Slug = *input.Slug
	}
	if input.Description != nil {
		c.Description = *input.Description
	}
	if input.IsActive != nil {
		c.IsActive = *input.IsActive
	}
	if input.SortOrder != nil {
		c.SortOrder = *input.SortOrder
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Category, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, activeOnly bool) ([]Category, error) {
	return s.repo.List(ctx, tenantID, activeOnly)
}
