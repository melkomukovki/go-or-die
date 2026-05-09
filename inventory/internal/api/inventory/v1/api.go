package v1

import inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer

	partService PartService
}

func NewAPI(partService PartService) *api {
	return &api{partService: partService}
}
