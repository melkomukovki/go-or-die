package order

import (
	"context"
	"fmt"
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

	switch order.Status {
	case model.OrderStatusPaid:
		return uuid.UUID{}, errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return uuid.UUID{}, errs.ErrOrderCancelled
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

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		return s.orderRepo.Update(ctx, order)
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("обновить заказ: %w", err)
	}

	return tId, nil
}
