package v1

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
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

func (c *client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
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
		OrderUuid:     orderUUID.String(),
		PaymentMethod: pMethod,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			return uuid.UUID{}, errs.ErrInvalidUUID
		}
		slog.Error("оплатить заказ", "error", err)
		return uuid.UUID{}, err
	}

	return uuid.MustParse(resp.TransactionUuid), nil
}
