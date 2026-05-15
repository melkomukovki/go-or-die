package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	"github.com/melkomukovki/go-or-die/order/internal/repository/converter"
	"github.com/melkomukovki/go-or-die/order/internal/repository/record"
)

type repository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func (r *repository) Create(ctx context.Context, order model.Order) error {
	orderRecord, itemRecord := converter.OrderToRecord(order)

	err := r.createOrder(ctx, orderRecord)
	if err != nil {
		return err
	}

	return r.createOrderItems(ctx, itemRecord)
}

func (r *repository) createOrder(ctx context.Context, order record.Order) error {
	query := squirrel.Insert("orders").
		Columns("uuid", "status", "transaction_uuid", "payment_method", "created_at").
		PlaceholderFormat(squirrel.Dollar).
		Values(order.UUID, order.Status, order.TransactionUUID, order.PaymentMethod, order.CreatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("подготовить запрос создания заказа: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}

func (r *repository) createOrderItems(ctx context.Context, items []record.OrderItem) error {
	query := squirrel.Insert("order_items").
		Columns("uuid", "order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range items {
		query = query.Values(item.UUID, item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("подготовить запрос создания предметов заказа: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("добавить предметы заказа: %w", err)
	}

	return nil
}

func (r *repository) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	queryOrder := squirrel.Select("uuid", "status", "transaction_uuid", "payment_method", "created_at", "updated_at").
		From("orders").Where(squirrel.Eq{"uuid": id}).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := queryOrder.ToSql()
	if err != nil {
		return model.Order{}, fmt.Errorf("подготовить запрос (заказ): %w", err)
	}

	var order record.Order
	err = r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, sql, args...).
		Scan(&order.UUID, &order.Status, &order.TransactionUUID, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, errs.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	queryOrderItems := squirrel.Select("uuid", "order_uuid", "part_uuid", "part_type", "price", "created_at").
		From("order_items").Where(squirrel.Eq{"order_uuid": order.UUID}).PlaceholderFormat(squirrel.Dollar)

	sql, args, err = queryOrderItems.ToSql()
	if err != nil {
		return model.Order{}, fmt.Errorf("подготовить запрос (детали заказа): %w", err)
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, sql, args...)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить детали заказа: %w", err)
	}
	defer rows.Close()

	orderItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return model.Order{}, fmt.Errorf("получить детали заказа: %w", err)
	}

	return converter.OrderToModel(order, orderItems), nil
}

func (r *repository) Update(ctx context.Context, order model.Order) error {
	query := squirrel.Update("orders").
		Set("status", order.Status).
		Set("transaction_uuid", order.TransactionUUID).
		Set("payment_method", order.PaymentMethod).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"uuid": order.UUID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("подготовить запрос (обновить заказ): %w", err)
	}
	tag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("обновить заказа: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}
