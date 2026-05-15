package order

import (
	"context"

	"github.com/melkomukovki/go-or-die/order/internal/api/order/v1/dto"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	order, err := a.orderService.Create(ctx, dto.OrderReqToModel(*req))
	if err != nil {
		return nil, err
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
	}, nil
}
