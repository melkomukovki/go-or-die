package dto

import (
	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/order/internal/model"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func OrderToDto(order model.Order) orderv1.OrderDto {
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(uuid.MustParse(*order.ShieldUUID))
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(uuid.MustParse(*order.WeaponUUID))
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(uuid.MustParse(*order.TransactionUUID))
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return orderv1.OrderDto{
		OrderUUID:       uuid.MustParse(order.UUID),
		HullUUID:        uuid.MustParse(order.HullUUID),
		EngineUUID:      uuid.MustParse(order.EngineUUID),
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderReqToModel(req orderv1.CreateOrderRequest) model.CreateOrderRequest {
	var shieldUUID *string
	if val, ok := req.ShieldUUID.Get(); ok && val != uuid.Nil {
		shieldUUID = new(val.String())
	}

	var weaponUUID *string
	if val, ok := req.WeaponUUID.Get(); ok && val != uuid.Nil {
		weaponUUID = new(val.String())
	}

	return model.CreateOrderRequest{
		HullUUID:   req.HullUUID.String(),
		EngineUUID: req.EngineUUID.String(),
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	}
}
