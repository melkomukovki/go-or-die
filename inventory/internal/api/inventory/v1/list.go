package v1

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/melkomukovki/go-or-die/inventory/internal/converter"
	errs "github.com/melkomukovki/go-or-die/inventory/internal/errors"
	"github.com/melkomukovki/go-or-die/inventory/internal/model"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	reqUuids := req.GetUuids()
	parsedUuids := make([]uuid.UUID, 0, len(reqUuids))
	if len(reqUuids) != 0 {
		for _, id := range reqUuids {
			parsed, err := uuid.Parse(id)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", id)
			}
			parsedUuids = append(parsedUuids, parsed)
		}
	}

	filter := model.PartFilter{
		UUIDs:    parsedUuids,
		PartType: converter.PartTypeToModel(req.GetPartType()),
	}

	parts, err := a.partService.List(ctx, filter)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrPartNotFound):
			return nil, status.Errorf(codes.NotFound, "деталь не найдена")
		default:
			slog.Error("получить деталь", "error", err)
			return nil, status.Error(codes.Internal, "внутренняя ошибка сервера")
		}
	}

	var partsProto []*inventoryv1.Part
	for _, part := range parts {
		partsProto = append(partsProto, converter.PartToProto(part))
	}
	return &inventoryv1.ListPartsResponse{Parts: partsProto}, nil
}
