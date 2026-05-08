package v1

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, param orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	transactionUuid, err := a.orderService.Pay(ctx, param.OrderUUID, model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderAlreadyPaid) || errors.Is(err, errs.ErrOrderPendingPaymentMismatch):
			return &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: errs.ErrOrderAlreadyPaid.Error(),
			}, nil
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.PayOrderNotFound{
				Code:    http.StatusNotFound,
				Message: errs.ErrOrderNotFound.Error(),
			}, nil
		case errors.Is(err, errs.ErrInvalidUUID):
			return &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: errs.ErrInvalidUUID.Error(),
			}, nil
		default:
			slog.Error("оплатить заказ", "error", err)
			return &orderv1.PayOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			}, nil
		}
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUuid,
	}, nil
}
