package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
)

func (s *service) Get(ctx context.Context, id string) (model.Order, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return model.Order{}, errs.ErrInvalidUUID
	}

	return s.orderRepo.Get(ctx, id)
}
