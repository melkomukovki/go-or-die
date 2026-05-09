package app

import (
	"google.golang.org/grpc"

	v1 "github.com/melkomukovki/go-or-die/payment/internal/api/payment/v1"
	"github.com/melkomukovki/go-or-die/payment/internal/service/payment"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

func Interceptors() []grpc.ServerOption {
	return nil
}

func RegisterServices(grpcServer *grpc.Server) {
	svc := payment.NewService()
	api := v1.NewAPI(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
}
