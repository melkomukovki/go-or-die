package v1

type api struct {
	orderService OrderService
}

func NewAPI(orderService OrderService) *api {
	return &api{orderService: orderService}
}
