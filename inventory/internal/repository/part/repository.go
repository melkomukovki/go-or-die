package part

import (
	"context"
	"sort"
	"time"

	errs "github.com/melkomukovki/go-or-die/inventory/internal/errors"
	"github.com/melkomukovki/go-or-die/inventory/internal/model"
	"github.com/melkomukovki/go-or-die/inventory/internal/repository/converter"
	"github.com/melkomukovki/go-or-die/inventory/internal/repository/record"
)

type repository struct {
	data map[string]record.Part
}

func NewRepository() *repository {
	now := time.Now()
	return &repository{
		data: map[string]record.Part{
			"550e8400-e29b-41d4-a716-446655440001": {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      "HULL",
				StockQuantity: 10,
				CreatedAt:     now,
			},
			"550e8400-e29b-41d4-a716-446655440002": {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      "HULL",
				StockQuantity: 5,
				CreatedAt:     now,
			},
			"550e8400-e29b-41d4-a716-446655440003": {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      "ENGINE",
				StockQuantity: 8,
				CreatedAt:     now,
			},
			"550e8400-e29b-41d4-a716-446655440004": {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      "ENGINE",
				StockQuantity: 3,
				CreatedAt:     now,
			},
			"550e8400-e29b-41d4-a716-446655440005": {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      "SHIELD",
				StockQuantity: 6,
				CreatedAt:     now,
			},
			"550e8400-e29b-41d4-a716-446655440006": {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      "WEAPON",
				StockQuantity: 7,
				CreatedAt:     now,
			},
			"550e8400-e29b-41d4-a716-446655440007": {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Прочный корпус для средних кораблей №2",
				Price:         2000000, // 20000₽
				PartType:      "HULL",
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

func (r *repository) Get(_ context.Context, uuid string) (model.Part, error) {
	if v, ok := r.data[uuid]; ok {
		return converter.PartToModel(v), nil
	}
	return model.Part{}, errs.ErrPartNotFound
}

func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) != 0 {
		var parts []model.Part
		for _, uuid := range filter.UUIDs {
			if v, ok := r.data[uuid]; ok {
				parts = append(parts, converter.PartToModel(v))
			} else {
				return nil, errs.ErrPartNotFound
			}
		}
		return parts, nil
	}

	partType := string(filter.PartType)
	var parts []model.Part
	for _, part := range r.data {
		if part.PartType == partType || filter.PartType == model.PartTypeUnspecified {
			parts = append(parts, converter.PartToModel(part))
		}
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return parts, nil
}
