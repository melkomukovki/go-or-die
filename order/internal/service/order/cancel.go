package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, id uuid.UUID) error {
	order, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return errs.ErrOrderPendingPaymentMismatch
	}

	order.Status = model.OrderStatusCancelled
	return s.orderRepo.Update(ctx, order)
}
