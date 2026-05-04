package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/go-or-die/payment/internal/converter"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	payRequest := converter.PayRequestToModel(req)

	transactionUuid, err := a.paymentService.Pay(ctx, payRequest)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &paymentv1.PayOrderResponse{TransactionUuid: transactionUuid}, nil
}
