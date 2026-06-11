package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusDraft    = "draft"
	StatusActive   = "active"
	StatusArchived = "archived"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_product_tenant_slug;index" json:"tenant_id"`
	CategoryID  *uuid.UUID     `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Slug        string         `gorm:"size:255;not null;uniqueIndex:idx_product_tenant_slug" json:"slug"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	SKU         string         `gorm:"size:100;index" json:"sku,omitempty"`
	Price       int64          `gorm:"not null;default:0" json:"price"`
	Status      string         `gorm:"size:50;not null;default:draft;index" json:"status"`
	IsFeatured  bool           `gorm:"not null;default:false" json:"is_featured"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Images    []ProductImage `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Inventory *Inventory     `gorm:"foreignKey:ProductID" json:"inventory,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Path      string    `gorm:"size:500;not null" json:"path"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (i *ProductImage) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type Inventory struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID          uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ProductID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"product_id"`
	Quantity          int       `gorm:"not null;default:0" json:"quantity"`
	LowStockThreshold int       `gorm:"not null;default:5" json:"low_stock_threshold"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (i *Inventory) IsLowStock() bool {
	return i.Quantity <= i.LowStockThreshold
}
