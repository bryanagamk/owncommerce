package database

import (
	"fmt"

	"github.com/agamtech/owncommerce/apps/api/internal/core/iam"
	"github.com/agamtech/owncommerce/apps/api/internal/core/subscription"
	"gorm.io/gorm"
)

var defaultPermissions = []iam.Permission{
	{Name: "product.view", Description: "View products", Module: "product"},
	{Name: "product.create", Description: "Create products", Module: "product"},
	{Name: "product.update", Description: "Update products", Module: "product"},
	{Name: "product.delete", Description: "Delete products", Module: "product"},
	{Name: "order.view", Description: "View orders", Module: "order"},
	{Name: "order.manage", Description: "Manage orders", Module: "order"},
	{Name: "customer.view", Description: "View customers", Module: "customer"},
	{Name: "subscription.manage", Description: "Manage subscription", Module: "subscription"},
	{Name: "tenant.manage", Description: "Manage tenant settings", Module: "tenant"},
	{Name: "staff.manage", Description: "Manage staff users", Module: "staff"},
	{Name: "audit.view", Description: "View audit logs", Module: "audit"},
	{Name: "platform.manage", Description: "Manage platform", Module: "platform"},
}

var defaultRoles = []struct {
	Name        string
	Scope       string
	Description string
	Permissions []string
}{
	{
		Name:        "super_admin",
		Scope:       "platform",
		Description: "Platform super administrator",
		Permissions: []string{"platform.manage", "audit.view", "subscription.manage", "tenant.manage"},
	},
	{
		Name:        "merchant_owner",
		Scope:       "tenant",
		Description: "Merchant store owner",
		Permissions: []string{
			"product.view", "product.create", "product.update", "product.delete",
			"order.view", "order.manage", "customer.view",
			"subscription.manage", "tenant.manage", "staff.manage", "audit.view",
		},
	},
	{
		Name:        "merchant_admin",
		Scope:       "tenant",
		Description: "Merchant administrator",
		Permissions: []string{
			"product.view", "product.create", "product.update", "product.delete",
			"order.view", "order.manage", "customer.view", "audit.view",
		},
	},
	{
		Name:        "merchant_staff",
		Scope:       "tenant",
		Description: "Merchant staff with limited access",
		Permissions: []string{"product.view", "order.view", "customer.view"},
	},
}

var defaultPlans = []struct {
	Name        string
	Slug        string
	Description string
	PriceMonthly int
	PriceYearly  int
	Features     []string
}{
	{
		Name:         "Trial",
		Slug:         "trial",
		Description:  "Free trial plan",
		PriceMonthly: 0,
		PriceYearly:  0,
		Features:     []string{"basic_storefront", "basic_analytics"},
	},
	{
		Name:         "Starter",
		Slug:         "starter",
		Description:  "Starter plan for small merchants",
		PriceMonthly: 99000,
		PriceYearly:  990000,
		Features:     []string{"basic_storefront", "basic_analytics"},
	},
	{
		Name:         "Growth",
		Slug:         "growth",
		Description:  "Growth plan with more features",
		PriceMonthly: 299000,
		PriceYearly:  2990000,
		Features:     []string{"basic_storefront", "basic_analytics", "custom_domain", "multi_staff"},
	},
	{
		Name:         "Enterprise",
		Slug:         "enterprise",
		Description:  "Enterprise plan",
		PriceMonthly: 999000,
		PriceYearly:  9990000,
		Features:     []string{"basic_storefront", "basic_analytics", "custom_domain", "multi_staff", "advanced_analytics", "marketplace_sync"},
	},
}

func Seed(db *gorm.DB) error {
	for _, perm := range defaultPermissions {
		if err := db.Where(iam.Permission{Name: perm.Name}).FirstOrCreate(&perm).Error; err != nil {
			return fmt.Errorf("seed permission %s: %w", perm.Name, err)
		}
	}

	for _, roleDef := range defaultRoles {
		role := iam.Role{
			Name:        roleDef.Name,
			Scope:       roleDef.Scope,
			Description: roleDef.Description,
			IsSystem:    true,
		}
		if err := db.Where(iam.Role{Name: role.Name, Scope: role.Scope}).FirstOrCreate(&role).Error; err != nil {
			return fmt.Errorf("seed role %s: %w", roleDef.Name, err)
		}

		for _, permName := range roleDef.Permissions {
			var perm iam.Permission
			if err := db.Where("name = ?", permName).First(&perm).Error; err != nil {
				return fmt.Errorf("find permission %s: %w", permName, err)
			}

			rp := iam.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
			if err := db.Where(rp).FirstOrCreate(&rp).Error; err != nil {
				return fmt.Errorf("seed role permission: %w", err)
			}
		}
	}

	for _, planDef := range defaultPlans {
		plan := subscription.Plan{
			Name:         planDef.Name,
			Slug:         planDef.Slug,
			Description:  planDef.Description,
			PriceMonthly: planDef.PriceMonthly,
			PriceYearly:  planDef.PriceYearly,
			IsActive:     true,
		}
		if err := db.Where(subscription.Plan{Slug: plan.Slug}).FirstOrCreate(&plan).Error; err != nil {
			return fmt.Errorf("seed plan %s: %w", planDef.Slug, err)
		}

		for _, feature := range planDef.Features {
			pf := subscription.PlanFeature{PlanID: plan.ID, FeatureKey: feature, Enabled: true}
			if err := db.Where(pf).FirstOrCreate(&pf).Error; err != nil {
				return fmt.Errorf("seed plan feature: %w", err)
			}
		}
	}

	return nil
}
