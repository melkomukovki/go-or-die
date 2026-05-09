package converter

import (
	"github.com/melkomukovki/go-or-die/payment/internal/model"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

func PayRequestToModel(req *paymentv1.PayOrderRequest) model.PayRequest {
	var paymentMethod model.PaymentMethod

	switch req.GetPaymentMethod() {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		paymentMethod = model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		paymentMethod = model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		paymentMethod = model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		paymentMethod = model.PaymentMethodInvestorMoney
	default:
		paymentMethod = model.PaymentMethodUnspecified
	}

	return model.PayRequest{OrderUUID: req.GetOrderUuid(), PaymentMethod: paymentMethod}
}
