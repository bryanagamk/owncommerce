package subscription

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plan struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Slug         string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Description  string    `gorm:"type:text" json:"description,omitempty"`
	PriceMonthly int       `gorm:"not null;default:0" json:"price_monthly"`
	PriceYearly  int       `gorm:"not null;default:0" json:"price_yearly"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Plan) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type PlanFeature struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PlanID     uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_plan_feature" json:"plan_id"`
	FeatureKey string    `gorm:"size:100;not null;uniqueIndex:idx_plan_feature" json:"feature_key"`
	Enabled    bool      `gorm:"not null;default:true" json:"enabled"`
	Value      string    `gorm:"size:255" json:"value,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (pf *PlanFeature) BeforeCreate(tx *gorm.DB) error {
	if pf.ID == uuid.Nil {
		pf.ID = uuid.New()
	}
	return nil
}

type Subscription struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID  uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"tenant_id"`
	PlanID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"plan_id"`
	Status    string     `gorm:"size:50;not null;default:trial" json:"status"`
	StartsAt  time.Time  `gorm:"not null" json:"starts_at"`
	EndsAt    *time.Time `json:"ends_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type FeatureFlag struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Key         string    `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Description string    `gorm:"size:500" json:"description,omitempty"`
	IsEnabled   bool      `gorm:"not null;default:false" json:"is_enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (f *FeatureFlag) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
