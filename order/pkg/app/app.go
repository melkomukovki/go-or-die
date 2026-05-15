package app

import (
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	apiv1 "github.com/melkomukovki/go-or-die/order/internal/api/order/v1"
	inventory "github.com/melkomukovki/go-or-die/order/internal/client/grpc/inventory/v1"
	payment "github.com/melkomukovki/go-or-die/order/internal/client/grpc/payment/v1"
	repository "github.com/melkomukovki/go-or-die/order/internal/repository/order"
	service "github.com/melkomukovki/go-or-die/order/internal/service/order"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(
	dbPool *pgxpool.Pool, txManager *manager.Manager, inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	repo := repository.NewRepository(dbPool)
	invClient := inventory.NewClientFromService(inventoryClient)
	payClient := payment.NewClientFromService(paymentClient)
	svc := service.NewService(invClient, payClient, repo, txManager)
	api := apiv1.NewAPI(svc)
	handler, err := orderv1.NewServer(api, orderv1.WithErrorHandler(apiv1.ErrorHandler))
	if err != nil {
		return nil, err
	}
	return handler, nil
}
