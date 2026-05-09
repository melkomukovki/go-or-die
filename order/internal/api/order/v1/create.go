package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/melkomukovki/go-or-die/order/internal/api/order/v1/dto"
	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	orderRequest := dto.OrderReqToModel(*req)
	order, err := a.orderService.Create(ctx, orderRequest)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidUUID):
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: errs.ErrInvalidUUID.Error(),
			}, nil
		case errors.Is(err, errs.ErrPartNotFound):
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: errs.ErrPartNotFound.Error(),
			}, nil
		case errors.Is(err, errs.ErrOutOfStock):
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: errs.ErrOutOfStock.Error(),
			}, nil
		default:
			slog.Error("создать заказ", "error", err)
			return &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			}, nil
		}
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
