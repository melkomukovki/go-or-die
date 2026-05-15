package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, req model.CreateOrderInput) (model.Order, error)
	Get(ctx context.Context, uuid uuid.UUID) (model.Order, error)
	Pay(ctx context.Context, uuid uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, uuid uuid.UUID) error
}
