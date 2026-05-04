package record

import "time"

type Order struct {
	OrderUUID       string
	HullUUID        string
	EngineUUID      string
	ShieldUUID      *string
	WeaponUUID      *string
	TotalPrice      int64
	TransactionUUID *string
	PaymentMethod   *string
	Status          *string
	CreatedAt       time.Time
}
