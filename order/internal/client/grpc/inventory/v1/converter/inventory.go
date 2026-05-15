package converter

import (
	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/order/internal/model"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
)

func PartsToModel(parts []*inventoryv1.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, model.Part{
			UUID:          uuid.MustParse(p.GetUuid()),
			Name:          p.GetName(),
			PartType:      partTypeFromProto(p.GetPartType()),
			Price:         p.GetPrice(),
			StockQuantity: p.GetStockQuantity(),
		})
	}
	return result
}

func partTypeFromProto(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}
