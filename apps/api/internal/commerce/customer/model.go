package customer

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID     uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_customer_tenant_email" json:"tenant_id"`
	Email        string         `gorm:"size:255;not null;uniqueIndex:idx_customer_tenant_email" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Phone        string         `gorm:"size:50" json:"phone,omitempty"`
	Status       string         `gorm:"size:50;not null;default:active" json:"status"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Address struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID      uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CustomerID    uuid.UUID `gorm:"type:uuid;not null;index" json:"customer_id"`
	Label         string    `gorm:"size:100" json:"label,omitempty"`
	RecipientName string    `gorm:"size:255;not null" json:"recipient_name"`
	Phone         string    `gorm:"size:50;not null" json:"phone"`
	AddressLine   string    `gorm:"type:text;not null" json:"address_line"`
	City          string    `gorm:"size:100;not null" json:"city"`
	Province      string    `gorm:"size:100;not null" json:"province"`
	PostalCode    string    `gorm:"size:20;not null" json:"postal_code"`
	IsDefault     bool      `gorm:"not null;default:false" json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

const StatusActive = "active"
