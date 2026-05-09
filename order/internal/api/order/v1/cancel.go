package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.CancelOrderNotFound{
				Code:    http.StatusNotFound,
				Message: errs.ErrOrderNotFound.Error(),
			}, nil
		case errors.Is(err, errs.ErrOrderPendingPaymentMismatch):
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: errs.ErrOrderPendingPaymentMismatch.Error(),
			}, nil
		default:
			slog.Error("отменить заказа", "error", err)
			return &orderv1.CancelOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			}, nil
		}
	}

	return &orderv1.CancelOrderResponse{}, nil
}
