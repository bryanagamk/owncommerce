package tenant

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Slug        string         `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Status      string         `gorm:"size:50;not null;default:active" json:"status"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	LogoURL     string         `gorm:"size:500" json:"logo_url,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TenantDomain struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID   uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Domain     string    `gorm:"size:255;uniqueIndex;not null" json:"domain"`
	Type       string    `gorm:"size:50;not null" json:"type"`
	IsPrimary  bool      `gorm:"not null;default:false" json:"is_primary"`
	IsVerified bool      `gorm:"not null;default:false" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (d *TenantDomain) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

const (
	DomainTypeSubdomain = "subdomain"
	DomainTypeCustom    = "custom"
	StatusActive        = "active"
	StatusSuspended     = "suspended"
)
