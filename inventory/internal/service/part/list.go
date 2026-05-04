package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/inventory/internal/errors"
	"github.com/melkomukovki/go-or-die/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) != 0 {
		for _, id := range filter.UUIDs {
			if _, err := uuid.Parse(id); err != nil {
				return nil, errs.ErrInvalidUUID
			}
		}
	}
	return s.partRepo.List(ctx, filter)
}
