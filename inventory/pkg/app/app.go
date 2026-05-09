package app

import (
	"google.golang.org/grpc"

	v1 "github.com/melkomukovki/go-or-die/inventory/internal/api/inventory/v1"
	repository "github.com/melkomukovki/go-or-die/inventory/internal/repository/part"
	service "github.com/melkomukovki/go-or-die/inventory/internal/service/part"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
)

func Interceptors() []grpc.ServerOption {
	return nil
}

func RegisterServices(grpcServer *grpc.Server) {
	repo := repository.NewRepository()
	svc := service.NewService(repo)
	api := v1.NewAPI(svc)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)
}
