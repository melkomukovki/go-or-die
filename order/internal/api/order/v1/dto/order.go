package dto

import (
	"github.com/google/uuid"

	"github.com/melkomukovki/go-or-die/order/internal/model"
	orderv1 "github.com/melkomukovki/go-or-die/shared/pkg/openapi/order/v1"
)

func OrderToDto(order model.Order) *orderv1.OrderDto {
	dto := &orderv1.OrderDto{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
		Status:     OrderStatusToDTO(order.Status),
		CreatedAt:  order.CreatedAt,
	}

	for _, item := range order.Items {
		switch item.PartType {
		case model.PartTypeHull:
			dto.HullUUID = item.PartUUID
		case model.PartTypeEngine:
			dto.EngineUUID = item.PartUUID
		case model.PartTypeShield:
			dto.ShieldUUID = orderv1.NewOptNilUUID(item.PartUUID)
		case model.PartTypeWeapon:
			dto.WeaponUUID = orderv1.NewOptNilUUID(item.PartUUID)
		}
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderv1.NewOptNilPaymentMethod(PaymentMethodToDTO(*order.PaymentMethod))
	}

	return dto
}

func OrderStatusToDTO(s model.OrderStatus) orderv1.OrderStatus {
	switch s {
	case model.OrderStatusPendingPayment:
		return orderv1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPaid:
		return orderv1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderv1.OrderStatusCANCELLED
	default:
		return orderv1.OrderStatus(s)
	}
}

func PaymentMethodToDTO(m model.PaymentMethod) orderv1.PaymentMethod {
	switch m {
	case model.PaymentMethodCard:
		return orderv1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderv1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderv1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderv1.PaymentMethodINVESTORMONEY
	default:
		return orderv1.PaymentMethod(m)
	}
}

func OrderReqToModel(req orderv1.CreateOrderRequest) model.CreateOrderInput {
	var shieldUUID *uuid.UUID
	if val, ok := req.ShieldUUID.Get(); ok && val != uuid.Nil {
		shieldUUID = &val
	}

	var weaponUUID *uuid.UUID
	if val, ok := req.WeaponUUID.Get(); ok && val != uuid.Nil {
		weaponUUID = &val
	}

	return model.CreateOrderInput{
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	}
}
