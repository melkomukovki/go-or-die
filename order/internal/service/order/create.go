package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
)

func (s *service) Create(ctx context.Context, req model.CreateOrderInput) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, req.PartUUIDs())
	if err != nil {
		if !errors.Is(err, errs.ErrPartNotFound) {
			slog.ErrorContext(ctx, "не удалось получить детали", "part_uuids", req.PartUUIDs(), "error", err)
		}
		return model.Order{}, fmt.Errorf("получить деталь: %w", err)
	}

	items := make([]model.OrderItem, 0, len(parts))
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrOutOfStock)
		}
		items = append(items, model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})
	}

	order := model.Order{
		UUID:      uuid.New(),
		Items:     items,
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		return s.orderRepo.Create(ctx, order)
	})
	if err != nil {
		return model.Order{}, fmt.Errorf("сохранить заказ: %w", err)
	}

	return order, nil
}
