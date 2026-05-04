package v1

import (
	"context"

	"github.com/melkomukovki/go-or-die/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error)
	Cancel(ctx context.Context, uuid string) error
}
