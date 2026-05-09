package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}
