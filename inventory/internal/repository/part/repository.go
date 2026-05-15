package part

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/melkomukovki/go-or-die/inventory/internal/errors"
	"github.com/melkomukovki/go-or-die/inventory/internal/model"
	"github.com/melkomukovki/go-or-die/inventory/internal/repository/converter"
	"github.com/melkomukovki/go-or-die/inventory/internal/repository/record"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) Get(ctx context.Context, id uuid.UUID) (model.Part, error) {
	query := `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts WHERE uuid = $1`

	var part record.Part
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&part.UUID, &part.Name, &part.Description, &part.PartType, &part.Price, &part.StockQuantity, &part.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Part{}, errs.ErrPartNotFound
		}
		return model.Part{}, err
	}

	return converter.PartToModel(part), nil
}

func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) != 0 {
		query := `
			SELECT p.uuid, p.name, p.description, p.part_type, p.price, p.stock_quantity, p.created_at
			FROM unnest($1::uuid[]) WITH ORDINALITY AS input(uuid, ord)
			JOIN parts p ON p.uuid = input.uuid
			ORDER BY input.ord
		`

		rows, err := r.pool.Query(ctx, query, filter.UUIDs)
		if err != nil {
			return []model.Part{}, err
		}
		defer rows.Close()

		parts := make([]model.Part, 0, len(filter.UUIDs))
		for rows.Next() {
			var part record.Part
			err = rows.Scan(&part.UUID, &part.Name, &part.Description, &part.PartType, &part.Price, &part.StockQuantity, &part.CreatedAt)
			if err != nil {
				return []model.Part{}, err
			}
			parts = append(parts, converter.PartToModel(part))
		}
		if err := rows.Err(); err != nil {
			return []model.Part{}, err
		}

		if len(parts) != len(filter.UUIDs) {
			return []model.Part{}, errs.ErrPartNotFound
		}
		return parts, nil
	}

	query := `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts`

	args := []any{}

	if filter.PartType != model.PartTypeUnspecified {
		query += ` WHERE part_type = $1`
		args = append(args, string(filter.PartType))
	}

	query += ` ORDER BY name`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return []model.Part{}, err
	}
	defer rows.Close()

	var parts []model.Part
	for rows.Next() {
		var part record.Part
		err = rows.Scan(
			&part.UUID, &part.Name, &part.Description, &part.PartType, &part.Price, &part.StockQuantity, &part.CreatedAt,
		)
		if err != nil {
			return []model.Part{}, err
		}
		parts = append(parts, converter.PartToModel(part))
	}

	if err = rows.Err(); err != nil {
		return []model.Part{}, err
	}

	return parts, nil
}
