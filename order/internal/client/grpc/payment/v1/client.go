package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/go-or-die/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

const grpcTimeout = time.Second * 5

type client struct {
	client paymentv1.PaymentServiceClient
}

func NewClientFromService(svc paymentv1.PaymentServiceClient) *client {
	return &client{client: svc}
}

func (c *client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	grpcCtx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := c.client.PayOrder(grpcCtx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: converter.PaymentMethodToDTO(method),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			return uuid.UUID{}, errs.ErrInvalidPaymentMethod
		}
		return uuid.UUID{}, fmt.Errorf("вызвать PaymentService.PayOrder: %w", err)
	}

	transactionUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("парсинг UUID транзакции из ответа PaymentService: %w", err)
	}

	return transactionUUID, nil
}
