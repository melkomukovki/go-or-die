package converter

import (
	"github.com/melkomukovki/go-or-die/order/internal/model"
	"github.com/melkomukovki/go-or-die/order/internal/repository/record"
)

func OrderToRecord(order model.Order) record.Order {
	var paymentMethod string
	if order.PaymentMethod == nil {
		paymentMethod = ""
	} else {
		paymentMethod = string(*order.PaymentMethod)
	}

	return record.Order{
		OrderUUID:       order.UUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          new(string(order.Status)),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderToModel(order record.Order) model.Order {
	return model.Order{
		UUID:            order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   new(model.PaymentMethod(*order.PaymentMethod)),
		Status:          model.OrderStatus(*order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
