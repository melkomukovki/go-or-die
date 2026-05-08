package v1

import (
	"context"

	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/go-or-die/payment/internal/converter"
	errs "github.com/melkomukovki/go-or-die/payment/internal/errors"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	payRequest := converter.PayRequestToModel(req)

	transactionUuid, err := a.paymentService.Pay(ctx, payRequest)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidPaymentMethod) || errors.Is(err, errs.ErrInvalidOrderUUID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &paymentv1.PayOrderResponse{TransactionUuid: transactionUuid}, nil
}
