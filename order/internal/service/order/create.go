package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
)

func (s *service) Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	uuids := []string{req.HullUUID, req.EngineUUID}
	if req.ShieldUUID != nil {
		uuids = append(uuids, *req.ShieldUUID)
	}
	if req.WeaponUUID != nil {
		uuids = append(uuids, *req.WeaponUUID)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	parts, err := s.inventoryClient.ListParts(ctx, uuids)
	if err != nil {
		return model.Order{}, err
	}

	var totalPrice int64
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, errs.ErrOutOfStock
		}
		totalPrice += part.Price
	}

	orderUuid := uuid.New()

	order := model.Order{
		UUID:       orderUuid.String(),
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: req.ShieldUUID,
		WeaponUUID: req.WeaponUUID,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
