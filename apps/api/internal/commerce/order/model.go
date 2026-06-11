package order

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusPendingPayment = "pending_payment"
	StatusPaid           = "paid"
	StatusProcessing     = "processing"
	StatusShipped        = "shipped"
	StatusCompleted      = "completed"
	StatusCancelled      = "cancelled"
)

type Order struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CustomerID      *uuid.UUID `gorm:"type:uuid;index" json:"customer_id,omitempty"`
	OrderNumber     string     `gorm:"size:50;not null;uniqueIndex" json:"order_number"`
	Status          string     `gorm:"size:50;not null;default:pending_payment;index" json:"status"`
	PaymentStatus   string     `gorm:"size:50;not null;default:pending" json:"payment_status"`
	Subtotal        int64      `gorm:"not null;default:0" json:"subtotal"`
	ShippingCost    int64      `gorm:"not null;default:0" json:"shipping_cost"`
	Total           int64      `gorm:"not null;default:0" json:"total"`
	RecipientName   string     `gorm:"size:255;not null" json:"recipient_name"`
	RecipientPhone  string     `gorm:"size:50;not null" json:"recipient_phone"`
	ShippingAddress string     `gorm:"type:text;not null" json:"shipping_address"`
	ShippingCity    string     `gorm:"size:100;not null" json:"shipping_city"`
	ShippingProvince string    `gorm:"size:100;not null" json:"shipping_province"`
	ShippingPostal  string     `gorm:"size:20;not null" json:"shipping_postal_code"`
	CustomerEmail   string     `gorm:"size:255" json:"customer_email,omitempty"`
	CustomerNote    string     `gorm:"type:text" json:"customer_note,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Items   []OrderItem          `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	History []OrderStatusHistory `gorm:"foreignKey:OrderID" json:"history,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type OrderItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	ProductName string    `gorm:"size:255;not null" json:"product_name"`
	SKU         string    `gorm:"size:100" json:"sku,omitempty"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	UnitPrice   int64     `gorm:"not null" json:"unit_price"`
	Subtotal    int64     `gorm:"not null" json:"subtotal"`
}

func (i *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type OrderStatusHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Status    string    `gorm:"size:50;not null" json:"status"`
	Note      string    `gorm:"type:text" json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *OrderStatusHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}
