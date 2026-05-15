package order

import (
	"context"

	"github.com/melkomukovki/go-or-die/order/internal/api/order/v1/dto"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		return nil, err
	}

	return dto.OrderToDto(order), nil
}
