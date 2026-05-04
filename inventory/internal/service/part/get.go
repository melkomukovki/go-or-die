package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/melkomukovki/go-or-die/inventory/internal/errors"
	"github.com/melkomukovki/go-or-die/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, id string) (model.Part, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}
	return s.partRepo.Get(ctx, id)
}
