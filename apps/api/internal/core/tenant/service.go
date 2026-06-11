package tenant

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

type Service struct {
	repo           *Repository
	platformDomain string
}

func NewService(repo *Repository, platformDomain string) *Service {
	return &Service{repo: repo, platformDomain: platformDomain}
}

type CreateTenantInput struct {
	Name string
	Slug string
}

func (s *Service) Create(ctx context.Context, input CreateTenantInput) (*Tenant, *TenantDomain, error) {
	slug := normalizeSlug(input.Slug)
	if slug == "" {
		return nil, nil, fmt.Errorf("slug is required")
	}
	if !slugPattern.MatchString(slug) {
		return nil, nil, fmt.Errorf("invalid slug format")
	}

	exists, err := s.repo.SlugExists(ctx, slug)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, fmt.Errorf("slug already taken")
	}

	domain := s.BuildSubdomain(slug)
	domainExists, err := s.repo.DomainExists(ctx, domain)
	if err != nil {
		return nil, nil, err
	}
	if domainExists {
		return nil, nil, fmt.Errorf("domain already taken")
	}

	tenant := &Tenant{
		Name:   input.Name,
		Slug:   slug,
		Status: StatusActive,
	}
	tenantDomain := &TenantDomain{
		Domain:     domain,
		Type:       DomainTypeSubdomain,
		IsPrimary:  true,
		IsVerified: true,
	}

	if err := s.repo.Create(ctx, tenant, tenantDomain); err != nil {
		return nil, nil, err
	}

	return tenant, tenantDomain, nil
}

func (s *Service) ResolveByHost(ctx context.Context, host string) (*Tenant, error) {
	host = normalizeHost(host)
	if host == "" {
		return nil, ErrTenantNotFound
	}
	return s.repo.FindByDomain(ctx, host)
}

func (s *Service) ResolveBySlug(ctx context.Context, slug string) (*Tenant, error) {
	return s.repo.FindBySlug(ctx, normalizeSlug(slug))
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	return s.repo.FindByID(ctx, id)
}

type UpdateStoreInput struct {
	Name         *string
	Description  *string
	ContactEmail *string
	ContactPhone *string
	Address      *string
	City         *string
	Province     *string
	PostalCode   *string
	LogoURL      *string
}

func (s *Service) UpdateStore(ctx context.Context, id uuid.UUID, input UpdateStoreInput) (*Tenant, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		t.Name = *input.Name
	}
	if input.Description != nil {
		t.Description = *input.Description
	}
	if input.ContactEmail != nil {
		t.ContactEmail = *input.ContactEmail
	}
	if input.ContactPhone != nil {
		t.ContactPhone = *input.ContactPhone
	}
	if input.Address != nil {
		t.Address = *input.Address
	}
	if input.City != nil {
		t.City = *input.City
	}
	if input.Province != nil {
		t.Province = *input.Province
	}
	if input.PostalCode != nil {
		t.PostalCode = *input.PostalCode
	}
	if input.LogoURL != nil {
		t.LogoURL = *input.LogoURL
	}
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) BuildSubdomain(slug string) string {
	return fmt.Sprintf("%s.%s", slug, s.platformDomain)
}

func normalizeSlug(slug string) string {
	return strings.ToLower(strings.TrimSpace(slug))
}

func normalizeHost(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return host
}
