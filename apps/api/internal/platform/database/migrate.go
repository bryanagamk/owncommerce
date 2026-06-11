package database

import (
	"fmt"

	"github.com/agamtech/owncommerce/apps/api/internal/core/audit"
	"github.com/agamtech/owncommerce/apps/api/internal/core/auth"
	"github.com/agamtech/owncommerce/apps/api/internal/core/iam"
	"github.com/agamtech/owncommerce/apps/api/internal/core/subscription"
	"github.com/agamtech/owncommerce/apps/api/internal/core/tenant"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&tenant.Tenant{},
		&tenant.TenantDomain{},
		&auth.User{},
		&auth.RefreshToken{},
		&iam.Role{},
		&iam.Permission{},
		&iam.RolePermission{},
		&iam.UserRole{},
		&audit.AuditLog{},
		&subscription.Plan{},
		&subscription.PlanFeature{},
		&subscription.Subscription{},
		&subscription.FeatureFlag{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}
