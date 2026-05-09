package payment

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/payment/internal/errors"
	"github.com/melkomukovki/go-or-die/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, req model.PayRequest) (string, error) {
	if req.OrderUUID == "" {
		return "", errs.ErrInvalidOrderUUID
	}
	if _, err := uuid.Parse(req.OrderUUID); err != nil {
		return "", errs.ErrInvalidOrderUUID
	}

	if !req.PaymentMethod.IsValid() {
		return "", errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New().String()

	slog.InfoContext(ctx, "оплата выполнена",
		"order_uuid", req.OrderUUID,
		"transaction_uuid", transactionUUID,
		"payment_method", req.PaymentMethod)

	return transactionUUID, nil
}
