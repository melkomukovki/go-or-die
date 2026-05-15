package order

type service struct {
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	orderRepo       OrderRepository
	txManager       TxManager
}

func NewService(inventoryClient InventoryClient, paymentClient PaymentClient, orderRepo OrderRepository, manager TxManager) *service {
	return &service{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepo:       orderRepo,
		txManager:       manager,
	}
}
