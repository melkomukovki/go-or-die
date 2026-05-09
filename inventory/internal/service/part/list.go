package part

import (
	"context"

	"github.com/melkomukovki/go-or-die/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	return s.partRepo.List(ctx, filter)
}
