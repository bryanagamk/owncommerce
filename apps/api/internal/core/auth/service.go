package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/agamtech/owncommerce/apps/api/internal/core/audit"
	"github.com/agamtech/owncommerce/apps/api/internal/core/iam"
	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
)

type Service struct {
	repo         *Repository
	jwt          *JWTManager
	tenantSvc    *tenant.Service
	iamRepo      *iam.Repository
	auditSvc     *audit.Service
}

func NewService(
	repo *Repository,
	jwt *JWTManager,
	tenantSvc *tenant.Service,
	iamRepo *iam.Repository,
	auditSvc *audit.Service,
) *Service {
	return &Service{
		repo:      repo,
		jwt:       jwt,
		tenantSvc: tenantSvc,
		iamRepo:   iamRepo,
		auditSvc:  auditSvc,
	}
}

type RegisterInput struct {
	Email     string
	Password  string
	Name      string
	StoreName string
	Slug      string
	UserAgent string
	IPAddress string
}

type AuthResult struct {
	User       *User
	Tenant     *tenant.Tenant
	Tokens     *TokenPair
	Permissions []string
	Roles      []string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	if _, err := s.repo.FindUserByEmail(ctx, input.Email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        input.Email,
		PasswordHash: string(hash),
		Name:         input.Name,
		Status:       UserStatusActive,
	}

	t, tenantDomain, err := s.tenantSvc.Create(ctx, tenant.CreateTenantInput{
		Name: input.StoreName,
		Slug: input.Slug,
	})
	if err != nil {
		return nil, err
	}
	_ = tenantDomain

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	role, err := s.iamRepo.FindRoleByName(ctx, "merchant_owner", "tenant")
	if err != nil {
		return nil, err
	}

	if err := s.iamRepo.AssignRole(ctx, user.ID, role.ID, &t.ID); err != nil {
		return nil, err
	}

	result, err := s.issueTokens(ctx, user, &t.ID, input.UserAgent, input.IPAddress)
	if err != nil {
		return nil, err
	}

	_ = s.auditSvc.Log(ctx, audit.LogInput{
		TenantID:   &t.ID,
		UserID:     &user.ID,
		Action:     "auth.register",
		EntityType: "user",
		EntityID:   user.ID.String(),
		Metadata: map[string]interface{}{
			"email": user.Email,
			"slug":  t.Slug,
		},
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
	})

	result.Tenant = t
	return result, nil
}

type LoginInput struct {
	Email     string
	Password  string
	UserAgent string
	IPAddress string
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	user, err := s.repo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tenantID, err := s.primaryTenantID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	result, err := s.issueTokens(ctx, user, tenantID, input.UserAgent, input.IPAddress)
	if err != nil {
		return nil, err
	}

	_ = s.repo.UpdateLastLogin(ctx, user.ID)
	_ = s.auditSvc.Log(ctx, audit.LogInput{
		TenantID:   tenantID,
		UserID:     &user.ID,
		Action:     "auth.login",
		EntityType: "user",
		EntityID:   user.ID.String(),
		IPAddress:  input.IPAddress,
		UserAgent:  input.UserAgent,
	})

	return result, nil
}

type RefreshInput struct {
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

func (s *Service) Refresh(ctx context.Context, input RefreshInput) (*AuthResult, error) {
	claims, err := s.jwt.ParseToken(input.RefreshToken)
	if err != nil || claims.TokenType != TokenTypeRefresh {
		return nil, fmt.Errorf("invalid refresh token")
	}

	stored, err := s.repo.FindRefreshTokenByHash(ctx, HashToken(input.RefreshToken))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}
	if stored.RevokedAt != nil || time.Now().UTC().After(stored.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired or revoked")
	}

	user, err := s.repo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return nil, err
	}

	result, err := s.issueTokens(ctx, user, claims.TenantID, input.UserAgent, input.IPAddress)
	if err != nil {
		return nil, err
	}

	_ = s.auditSvc.Log(ctx, audit.LogInput{
		TenantID:   claims.TenantID,
		UserID:     &user.ID,
		Action:     "auth.refresh",
		EntityType: "user",
		EntityID:   user.ID.String(),
		IPAddress:  input.IPAddress,
		UserAgent:  input.UserAgent,
	})

	return result, nil
}

func (s *Service) Logout(ctx context.Context, userID uuid.UUID, refreshToken string, ip, ua string) error {
	if refreshToken != "" {
		if stored, err := s.repo.FindRefreshTokenByHash(ctx, HashToken(refreshToken)); err == nil {
			_ = s.repo.RevokeRefreshToken(ctx, stored.ID)
		}
	} else {
		_ = s.repo.RevokeAllUserTokens(ctx, userID)
	}

	_ = s.auditSvc.Log(ctx, audit.LogInput{
		UserID:     &userID,
		Action:     "auth.logout",
		EntityType: "user",
		EntityID:   userID.String(),
		IPAddress:  ip,
		UserAgent:  ua,
	})

	return nil
}

func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	return s.repo.FindUserByID(ctx, userID)
}

func (s *Service) issueTokens(ctx context.Context, user *User, tenantID *uuid.UUID, userAgent, ip string) (*AuthResult, error) {
	tokens, refreshHash, err := s.jwt.GenerateTokenPair(user.ID, tenantID)
	if err != nil {
		return nil, err
	}

	refresh := &RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().UTC().Add(s.jwt.RefreshExpiry()),
		UserAgent: userAgent,
		IPAddress: ip,
	}
	if err := s.repo.CreateRefreshToken(ctx, refresh); err != nil {
		return nil, err
	}

	perms, roles, err := s.iamRepo.GetUserPermissions(ctx, user.ID, tenantID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:        user,
		Tokens:      tokens,
		Permissions: perms,
		Roles:       roles,
	}, nil
}

func (s *Service) primaryTenantID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	tenantID, err := s.iamRepo.GetPrimaryTenantID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tenantID, nil
}
