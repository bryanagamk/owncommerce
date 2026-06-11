package audit

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type LogInput struct {
	TenantID   *uuid.UUID
	UserID     *uuid.UUID
	Action     string
	EntityType string
	EntityID   string
	Metadata   map[string]interface{}
	IPAddress  string
	UserAgent  string
}

func (s *Service) Log(ctx context.Context, input LogInput) error {
	var metadata datatypes.JSON
	if input.Metadata != nil {
		raw, err := json.Marshal(input.Metadata)
		if err != nil {
			return err
		}
		metadata = datatypes.JSON(raw)
	}

	log := AuditLog{
		TenantID:   input.TenantID,
		UserID:     input.UserID,
		Action:     input.Action,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		Metadata:   metadata,
		IPAddress:  input.IPAddress,
		UserAgent:  input.UserAgent,
	}

	return s.db.WithContext(ctx).Create(&log).Error
}

type ListFilter struct {
	TenantID *uuid.UUID
	UserID   *uuid.UUID
	Action   string
	Limit    int
	Offset   int
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]AuditLog, int64, error) {
	query := s.db.WithContext(ctx).Model(&AuditLog{})
	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var logs []AuditLog
	if err := query.Order("created_at DESC").Limit(limit).Offset(filter.Offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
