package order

type service struct {
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	orderRepo       OrderRepository
}

func NewService(inventoryClient InventoryClient, paymentClient PaymentClient, orderRepo OrderRepository) *service {
	return &service{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepo:       orderRepo,
	}
}
