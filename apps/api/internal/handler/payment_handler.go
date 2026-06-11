package handler

import (
	"github.com/agamtech/owncommerce/apps/api/internal/commerce/payment"
	"github.com/agamtech/owncommerce/apps/api/internal/platform/response"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentSvc *payment.Service
}

func NewPaymentHandler(paymentSvc *payment.Service) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

func (h *PaymentHandler) MidtransWebhook(c *fiber.Ctx) error {
	var payload payment.NotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.BadRequest(c, "invalid notification payload")
	}

	if payload.SignatureKey != "" && !h.paymentSvc.VerifySignature(payload) {
		return response.Forbidden(c, "invalid signature")
	}

	if err := h.paymentSvc.HandleNotification(c.UserContext(), payload); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{"message": "notification processed"})
}
