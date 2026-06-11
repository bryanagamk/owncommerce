package iam

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrRoleNotFound = errors.New("role not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindRoleByName(ctx context.Context, name, scope string) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Where("name = ? AND scope = ?", name, scope).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *Repository) AssignRole(ctx context.Context, userID, roleID uuid.UUID, tenantID *uuid.UUID) error {
	ur := UserRole{
		UserID:   userID,
		RoleID:   roleID,
		TenantID: tenantID,
	}
	return r.db.WithContext(ctx).Where(ur).FirstOrCreate(&ur).Error
}

func (r *Repository) GetPrimaryTenantID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	var ur UserRole
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND tenant_id IS NOT NULL", userID).
		Order("created_at ASC").
		First(&ur).Error
	if err != nil {
		return nil, err
	}
	return ur.TenantID, nil
}

func (r *Repository) GetUserPermissions(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID) ([]string, []string, error) {
	var userRoles []UserRole
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if tenantID != nil {
		query = query.Where("tenant_id IS NULL OR tenant_id = ?", *tenantID)
	}
	if err := query.Find(&userRoles).Error; err != nil {
		return nil, nil, err
	}

	if len(userRoles) == 0 {
		return []string{}, []string{}, nil
	}

	roleIDs := make([]uuid.UUID, 0, len(userRoles))
	roleMap := make(map[uuid.UUID]struct{})
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
		roleMap[ur.RoleID] = struct{}{}
	}

	var roles []Role
	if err := r.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, nil, err
	}

	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	var rolePerms []RolePermission
	if err := r.db.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&rolePerms).Error; err != nil {
		return nil, nil, err
	}

	permIDs := make([]uuid.UUID, 0, len(rolePerms))
	for _, rp := range rolePerms {
		permIDs = append(permIDs, rp.PermissionID)
	}

	if len(permIDs) == 0 {
		return []string{}, roleNames, nil
	}

	var perms []Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", permIDs).Find(&perms).Error; err != nil {
		return nil, nil, err
	}

	permNames := make([]string, 0, len(perms))
	for _, p := range perms {
		permNames = append(permNames, p.Name)
	}

	return permNames, roleNames, nil
}

func (r *Repository) UserHasPermission(ctx context.Context, userID uuid.UUID, tenantID *uuid.UUID, permission string) (bool, error) {
	perms, _, err := r.GetUserPermissions(ctx, userID, tenantID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		if p == permission {
			return true, nil
		}
	}
	return false, nil
}

func (r *Repository) ListRoles(ctx context.Context) ([]Role, error) {
	var roles []Role
	if err := r.db.WithContext(ctx).Order("scope, name").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) ListPermissions(ctx context.Context) ([]Permission, error) {
	var perms []Permission
	if err := r.db.WithContext(ctx).Order("module, name").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}
