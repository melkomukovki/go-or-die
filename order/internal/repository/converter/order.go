package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/order/internal/model"
	"github.com/melkomukovki/go-or-die/order/internal/repository/record"
)

func OrderToRecord(order model.Order) (record.Order, []record.OrderItem) {
	orderRecord := record.Order{
		UUID:      order.UUID,
		Status:    string(order.Status),
		CreatedAt: order.CreatedAt,
	}

	if order.PaymentMethod != nil {
		orderRecord.PaymentMethod = new(string(*order.PaymentMethod))
	}

	if order.PaymentMethod != nil {
		orderRecord.PaymentMethod = new(string(*order.PaymentMethod))
	}

	now := time.Now()

	items := make([]record.OrderItem, 0, len(order.Items))
	for _, i := range order.Items {
		item := record.OrderItem{
			UUID:      uuid.New(),
			OrderUUID: order.UUID,
			PartUUID:  i.PartUUID,
			PartType:  string(i.PartType),
			Price:     i.Price,
			CreatedAt: now,
		}
		items = append(items, item)
	}

	return orderRecord, items
}

func OrderToModel(or record.Order, oir []record.OrderItem) model.Order {
	order := model.Order{
		UUID:            or.UUID,
		TransactionUUID: or.TransactionUUID,
		Status:          model.OrderStatus(or.Status),
		CreatedAt:       or.CreatedAt,
	}

	if or.PaymentMethod != nil {
		order.PaymentMethod = new(model.PaymentMethod(*or.PaymentMethod))
	}

	items := make([]model.OrderItem, 0, len(oir))
	for _, i := range oir {
		item := model.OrderItem{
			PartUUID: i.PartUUID,
			PartType: model.PartType(i.PartType),
			Price:    i.Price,
		}
		items = append(items, item)
	}

	order.Items = items
	return order
}
