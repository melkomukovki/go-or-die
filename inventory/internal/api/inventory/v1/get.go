package v1

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/go-or-die/inventory/internal/converter"
	errs "github.com/melkomukovki/go-or-die/inventory/internal/errors"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	paramUuid := req.GetUuid()
	part, err := a.partService.Get(ctx, paramUuid)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidUUID):
			return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", paramUuid)
		case errors.Is(err, errs.ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "деталь c uuid=%s не найдена", paramUuid)
		default:
			slog.Error("получить деталь", "error", err)
			return nil, status.Error(codes.Internal, "внутренняя ошибка сервера")
		}
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
