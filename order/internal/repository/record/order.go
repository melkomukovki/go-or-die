package record

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID
	WeaponUUID      *uuid.UUID
	TotalPrice      int64
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          *string
	CreatedAt       time.Time
}
