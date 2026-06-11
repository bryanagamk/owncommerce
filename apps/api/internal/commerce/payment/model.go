package payment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	StatusPending = "pending"
	StatusPaid    = "paid"
	StatusFailed  = "failed"
	StatusExpired = "expired"
)

type Payment struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"tenant_id"`
	OrderID           uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	MidtransOrderID   string         `gorm:"size:100;not null;uniqueIndex" json:"midtrans_order_id"`
	SnapToken         string         `gorm:"size:500" json:"snap_token,omitempty"`
	Status            string         `gorm:"size:50;not null;default:pending;index" json:"status"`
	TransactionStatus string         `gorm:"size:50" json:"transaction_status,omitempty"`
	PaymentType       string         `gorm:"size:50" json:"payment_type,omitempty"`
	GrossAmount       int64          `gorm:"not null;default:0" json:"gross_amount"`
	RawNotification   datatypes.JSON `gorm:"type:jsonb" json:"raw_notification,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
