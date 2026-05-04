package v1

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

type client struct {
	client paymentv1.PaymentServiceClient
}

func NewClientFromService(svc paymentv1.PaymentServiceClient) *client {
	return &client{client: svc}
}

func (c *client) PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	var pMethod paymentv1.PaymentMethod
	switch method {
	case model.PaymentMethodCard:
		pMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		pMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		pMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		pMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		pMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	resp, err := c.client.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: pMethod,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			return "", errs.ErrInvalidUUID
		}
		slog.Error("оплатить заказ", "error", err)
		return "", err
	}

	return resp.TransactionUuid, nil
}
