package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
)

func (s *service) Pay(ctx context.Context, id uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return uuid.UUID{}, err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return uuid.UUID{}, errs.ErrOrderPendingPaymentMismatch
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	tId, err := s.paymentClient.PayOrder(ctx, id, method)
	if err != nil {
		return uuid.UUID{}, err
	}

	order.Status = model.OrderStatusPaid
	order.PaymentMethod = &method
	order.TransactionUUID = &tId

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return uuid.UUID{}, err
	}

	return tId, nil
}
