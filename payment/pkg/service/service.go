package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

// PaymentServer реализует gRPC сервис оплаты.
type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
}

// PayOrder обрабатывает оплату заказа.
func (s *PaymentServer) PayOrder(
	_ context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	orderUuid := req.GetOrderUuid()
	if orderUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid не может быть пустым")
	}

	if _, err := uuid.Parse(orderUuid); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "некорректное значение uuid=%s", orderUuid)
	}

	paymentMethod := req.GetPaymentMethod()
	if paymentMethod == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "не указан метод оплаты")
	}

	transactionUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при формировании UUIDv4: %s", err.Error())
	}

	slog.Info("оплата прошла успешно",
		"order_uuid", req.GetOrderUuid(),
		"transaction_uuid", transactionUuid.String(),
	)

	return &paymentv1.PayOrderResponse{TransactionUuid: transactionUuid.String()}, nil
}
