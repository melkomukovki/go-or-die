package v1

import (
	"context"

	"github.com/melkomukovki/go-or-die/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, req model.PayRequest) (string, error)
}
