package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/order/internal/model"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	return s.orderRepo.Get(ctx, id)
}
