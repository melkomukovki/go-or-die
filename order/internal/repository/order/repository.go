package order

import (
	"context"
	"sync"

	errs "github.com/melkomukovki/go-or-die/order/internal/errors"
	"github.com/melkomukovki/go-or-die/order/internal/model"
	"github.com/melkomukovki/go-or-die/order/internal/repository/converter"
	"github.com/melkomukovki/go-or-die/order/internal/repository/record"
)

type repository struct {
	mu     sync.RWMutex
	orders map[string]record.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[string]record.Order),
	}
}

func (r *repository) Create(ctx context.Context, order model.Order) error {
	orderRecord := converter.OrderToRecord(order)
	r.mu.Lock()
	r.orders[order.UUID] = orderRecord
	r.mu.Unlock()
	return nil
}

func (r *repository) Get(ctx context.Context, uuid string) (model.Order, error) {
	r.mu.RLock()
	orderRecord, ok := r.orders[uuid]
	r.mu.RUnlock()

	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return converter.OrderToModel(orderRecord), nil
}

func (r *repository) Update(ctx context.Context, order model.Order) error {
	r.mu.RLock()
	_, ok := r.orders[order.UUID]
	r.mu.RUnlock()

	if !ok {
		return errs.ErrOrderNotFound
	}

	r.mu.Lock()
	r.orders[order.UUID] = converter.OrderToRecord(order)
	r.mu.Unlock()

	return nil
}
