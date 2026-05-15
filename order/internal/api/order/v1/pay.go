package order

import (
	"context"

	"github.com/melkomukovki/go-or-die/order/internal/model"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, param orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	transactionUuid, err := a.orderService.Pay(ctx, param.OrderUUID, model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		return nil, err
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUuid,
	}, nil
}
