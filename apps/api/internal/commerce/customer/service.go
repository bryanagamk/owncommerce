package customer

import (
	"context"
	"errors"
	"fmt"

	"github.com/agamtech/owncommerce/apps/api/internal/core/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service struct {
	repo *Repository
	jwt  *auth.JWTManager
}

func NewService(repo *Repository, jwt *auth.JWTManager) *Service {
	return &Service{repo: repo, jwt: jwt}
}

type RegisterInput struct {
	TenantID uuid.UUID
	Email    string
	Password string
	Name     string
	Phone    string
}

type AuthResult struct {
	Customer *Customer
	Tokens   *auth.TokenPair
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	if _, err := s.repo.FindByEmail(ctx, input.TenantID, input.Email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	c := &Customer{
		TenantID:     input.TenantID,
		Email:        input.Email,
		PasswordHash: string(hash),
		Name:         input.Name,
		Phone:        input.Phone,
		Status:       StatusActive,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	tokens, err := s.jwt.GenerateCustomerTokenPair(c.ID, c.TenantID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Customer: c, Tokens: tokens}, nil
}

type LoginInput struct {
	TenantID uuid.UUID
	Email    string
	Password string
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	c, err := s.repo.FindByEmail(ctx, input.TenantID, input.Email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.jwt.GenerateCustomerTokenPair(c.ID, c.TenantID)
	if err != nil {
		return nil, err
	}
	_ = s.repo.UpdateLastLogin(ctx, c.ID)

	return &AuthResult{Customer: c, Tokens: tokens}, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*Customer, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

type UpdateProfileInput struct {
	Name  *string
	Phone *string
}

func (s *Service) UpdateProfile(ctx context.Context, tenantID, id uuid.UUID, input UpdateProfileInput) (*Customer, error) {
	c, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		c.Name = *input.Name
	}
	if input.Phone != nil {
		c.Phone = *input.Phone
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

type AddressInput struct {
	Label         string `json:"label"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	AddressLine   string `json:"address_line"`
	City          string `json:"city"`
	Province      string `json:"province"`
	PostalCode    string `json:"postal_code"`
	IsDefault     bool   `json:"is_default"`
}

func (s *Service) CreateAddress(ctx context.Context, tenantID, customerID uuid.UUID, input AddressInput) (*Address, error) {
	a := &Address{
		TenantID:      tenantID,
		CustomerID:    customerID,
		Label:         input.Label,
		RecipientName: input.RecipientName,
		Phone:         input.Phone,
		AddressLine:   input.AddressLine,
		City:          input.City,
		Province:      input.Province,
		PostalCode:    input.PostalCode,
		IsDefault:     input.IsDefault,
	}
	if err := s.repo.CreateAddress(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) UpdateAddress(ctx context.Context, tenantID, customerID, id uuid.UUID, input AddressInput) (*Address, error) {
	a, err := s.repo.FindAddress(ctx, tenantID, customerID, id)
	if err != nil {
		return nil, err
	}
	a.Label = input.Label
	a.RecipientName = input.RecipientName
	a.Phone = input.Phone
	a.AddressLine = input.AddressLine
	a.City = input.City
	a.Province = input.Province
	a.PostalCode = input.PostalCode
	a.IsDefault = input.IsDefault
	if err := s.repo.UpdateAddress(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAddress(ctx context.Context, tenantID, customerID, id uuid.UUID) error {
	return s.repo.DeleteAddress(ctx, tenantID, customerID, id)
}

func (s *Service) ListAddresses(ctx context.Context, tenantID, customerID uuid.UUID) ([]Address, error) {
	return s.repo.ListAddresses(ctx, tenantID, customerID)
}

func (s *Service) GetAddress(ctx context.Context, tenantID, customerID, id uuid.UUID) (*Address, error) {
	a, err := s.repo.FindAddress(ctx, tenantID, customerID, id)
	if err != nil {
		return nil, err
	}
	if a.CustomerID != customerID {
		return nil, fmt.Errorf("address not found")
	}
	return a, nil
}
