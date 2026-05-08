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

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.GetOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "заказ не найден",
			}, nil
		default:
			slog.Error("получить информацию о заказе", "error", err)
			return &orderv1.GetOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			}, nil
		}
	}

	return new(dto.OrderToDto(order)), nil
}
