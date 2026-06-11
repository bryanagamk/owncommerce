package payment

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/agamtech/owncommerce/apps/api/internal/commerce/order"
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/product"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Service struct {
	repo        *Repository
	orderRepo   *order.Repository
	productRepo *product.Repository
	midtrans    *MidtransClient
}

func NewService(repo *Repository, orderRepo *order.Repository, productRepo *product.Repository, midtrans *MidtransClient) *Service {
	return &Service{repo: repo, orderRepo: orderRepo, productRepo: productRepo, midtrans: midtrans}
}

type PayResult struct {
	SnapToken   string `json:"snap_token"`
	RedirectURL string `json:"redirect_url,omitempty"`
	PaymentID   uuid.UUID `json:"payment_id"`
	OrderID     uuid.UUID `json:"order_id"`
}

func (s *Service) CreateSnapPayment(ctx context.Context, tenantID, orderID uuid.UUID) (*PayResult, error) {
	if s.midtrans == nil || s.midtrans.cfg.ServerKey == "" {
		return nil, fmt.Errorf("midtrans is not configured")
	}

	o, err := s.orderRepo.FindByID(ctx, tenantID, orderID)
	if err != nil {
		return nil, err
	}
	if o.Status != order.StatusPendingPayment {
		return nil, fmt.Errorf("order is not pending payment")
	}

	midtransOrderID := o.OrderNumber
	if existing, err := s.repo.FindByOrderID(ctx, tenantID, orderID); err == nil {
		if existing.SnapToken != "" && existing.Status == StatusPending {
			return &PayResult{
				SnapToken: existing.SnapToken,
				PaymentID: existing.ID,
				OrderID:   orderID,
			}, nil
		}
		midtransOrderID = existing.MidtransOrderID
	}

	items := make([]SnapItemDetail, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, SnapItemDetail{
			ID:       item.ProductID.String(),
			Price:    item.UnitPrice,
			Quantity: item.Quantity,
			Name:     item.ProductName,
		})
	}
	if o.ShippingCost > 0 {
		items = append(items, SnapItemDetail{
			ID:       "shipping",
			Price:    o.ShippingCost,
			Quantity: 1,
			Name:     "Shipping",
		})
	}

	snapReq := SnapTransactionRequest{
		TransactionDetails: SnapTransactionDetails{
			OrderID:     midtransOrderID,
			GrossAmount: o.Total,
		},
		CustomerDetails: &SnapCustomerDetails{
			FirstName: o.RecipientName,
			Email:     o.CustomerEmail,
			Phone:     o.RecipientPhone,
		},
		ItemDetails: items,
	}

	snapResp, err := s.midtrans.CreateSnapToken(snapReq)
	if err != nil {
		return nil, err
	}

	p := &Payment{
		TenantID:        tenantID,
		OrderID:         orderID,
		MidtransOrderID: midtransOrderID,
		SnapToken:       snapResp.Token,
		Status:          StatusPending,
		GrossAmount:     o.Total,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		// upsert if exists
		if existing, findErr := s.repo.FindByOrderID(ctx, tenantID, orderID); findErr == nil {
			existing.SnapToken = snapResp.Token
			existing.MidtransOrderID = midtransOrderID
			existing.GrossAmount = o.Total
			if err := s.repo.Update(ctx, existing); err != nil {
				return nil, err
			}
			p = existing
		} else {
			return nil, err
		}
	}

	return &PayResult{
		SnapToken:   snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
		PaymentID:   p.ID,
		OrderID:     orderID,
	}, nil
}

func (s *Service) HandleNotification(ctx context.Context, payload NotificationPayload) error {
	if s.midtrans == nil {
		return fmt.Errorf("midtrans is not configured")
	}

	p, err := s.repo.FindByMidtransOrderID(ctx, payload.OrderID)
	if err != nil {
		return err
	}

	raw, _ := json.Marshal(payload)
	p.RawNotification = datatypes.JSON(raw)
	p.TransactionStatus = payload.TransactionStatus
	p.PaymentType = payload.PaymentType

	switch payload.TransactionStatus {
	case "capture", "settlement":
		p.Status = StatusPaid
		if err := s.repo.Update(ctx, p); err != nil {
			return err
		}
		if err := s.orderRepo.MarkPaid(ctx, p.TenantID, p.OrderID); err != nil {
			return err
		}
		o, err := s.orderRepo.FindByID(ctx, p.TenantID, p.OrderID)
		if err != nil {
			return err
		}
		for _, item := range o.Items {
			_ = s.productRepo.DecrementStock(ctx, p.TenantID, item.ProductID, item.Quantity)
		}
	case "deny", "cancel", "expire":
		p.Status = StatusFailed
		_ = s.repo.Update(ctx, p)
		_, _ = s.orderRepo.UpdateStatus(ctx, p.TenantID, p.OrderID, order.StatusCancelled, "payment "+payload.TransactionStatus)
	default:
		_ = s.repo.Update(ctx, p)
	}

	return nil
}

func (s *Service) VerifySignature(payload NotificationPayload) bool {
	gross, _ := strconv.ParseFloat(payload.GrossAmount, 64)
	grossStr := fmt.Sprintf("%.2f", gross)
	if payload.GrossAmount == strconv.FormatInt(int64(gross), 10) {
		grossStr = payload.GrossAmount
	}

	raw := payload.OrderID + payload.StatusCode + grossStr + s.midtrans.cfg.ServerKey
	sum := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(sum[:])
	return expected == payload.SignatureKey
}
