package order

import (
	"context"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return errs.ErrOrderPendingPaymentMismatch
	}

	order.Status = model.OrderStatusCancelled
	return s.orderRepo.Update(ctx, order)
}
