package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/go-or-die/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
)

const grpcTimeout = time.Second * 5

type client struct {
	client inventoryv1.InventoryServiceClient
}

func NewClientFromService(svc inventoryv1.InventoryServiceClient) *client {
	return &client{client: svc}
}

func (c *client) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	uuidsStr := make([]string, len(uuids))
	for i, id := range uuids {
		uuidsStr[i] = id.String()
	}

	grpcCtx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := c.client.ListParts(grpcCtx, &inventoryv1.ListPartsRequest{
		Uuids: uuidsStr,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, errs.ErrPartNotFound
		}
		return nil, fmt.Errorf("вызвать InventorySerive.ListParts: %w", err)
	}

	return converter.PartsToModel(resp.GetParts()), nil
}
