package audit

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID   *uuid.UUID     `gorm:"type:uuid;index" json:"tenant_id,omitempty"`
	UserID     *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Action     string         `gorm:"size:100;not null;index" json:"action"`
	EntityType string         `gorm:"size:100" json:"entity_type,omitempty"`
	EntityID   string         `gorm:"size:100" json:"entity_id,omitempty"`
	Metadata   datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	IPAddress  string         `gorm:"size:100" json:"ip_address,omitempty"`
	UserAgent  string         `gorm:"size:500" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `gorm:"index" json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
