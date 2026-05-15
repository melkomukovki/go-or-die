package order

import orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"

type api struct {
	orderv1.UnimplementedHandler
	orderService OrderService
}

func NewAPI(orderService OrderService) *api {
	return &api{orderService: orderService}
}
