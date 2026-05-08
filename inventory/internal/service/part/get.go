package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (model.Part, error) {
	return s.partRepo.Get(ctx, id)
}
