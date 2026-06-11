package cart

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const StatusActive = "active"

type Cart struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CustomerID *uuid.UUID `gorm:"type:uuid;index" json:"customer_id,omitempty"`
	SessionID  string     `gorm:"size:100;index" json:"session_id,omitempty"`
	Status     string     `gorm:"size:50;not null;default:active" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Items []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_cart_product" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_cart_product" json:"product_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	UnitPrice int64     `gorm:"not null" json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i *CartItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
