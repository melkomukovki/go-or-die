package v1

import paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"

type api struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService PaymentService
}

func NewAPI(paymentService PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
