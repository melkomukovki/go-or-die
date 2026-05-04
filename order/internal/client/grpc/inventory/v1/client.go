package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
)

type client struct {
	client inventoryv1.InventoryServiceClient
}

func NewClientFromService(svc inventoryv1.InventoryServiceClient) *client {
	return &client{client: svc}
}

func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	resp, err := c.client.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, errs.ErrPartNotFound
		}
		return nil, err
	}
	parts := resp.GetParts()
	var partsResp []model.Part
	for _, part := range parts {
		modelPart := model.Part{
			UUID:          part.Uuid,
			Name:          part.Name,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
		}
		partsResp = append(partsResp, modelPart)
	}
	return partsResp, nil
}
