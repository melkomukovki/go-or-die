package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID            uuid.UUID
	Items           []OrderItem
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

func (o Order) TotalPrice() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price
	}
	return total
}

type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type CreateOrderInput struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

func (i *CreateOrderInput) PartUUIDs() []uuid.UUID {
	uuids := []uuid.UUID{i.HullUUID, i.EngineUUID}
	if i.ShieldUUID != nil {
		uuids = append(uuids, *i.ShieldUUID)
	}
	if i.WeaponUUID != nil {
		uuids = append(uuids, *i.WeaponUUID)
	}
	return uuids
}
