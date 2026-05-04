package converter

import (
	"github.com/melkomukovki/go-or-die/inventory/internal/model"
	"github.com/melkomukovki/go-or-die/inventory/internal/repository/record"
)

func PartToModel(part record.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}
