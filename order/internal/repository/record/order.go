package record

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID            uuid.UUID  `db:"uuid"`
	Status          string     `db:"status"`
	TransactionUUID *uuid.UUID `db:"transaction_uuid"`
	PaymentMethod   *string    `db:"payment_method"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

type OrderItem struct {
	UUID      uuid.UUID `db:"uuid"`
	OrderUUID uuid.UUID `db:"order_uuid"`
	PartUUID  uuid.UUID `db:"part_uuid"`
	PartType  string    `db:"part_type"`
	Price     int64     `db:"price"`
	CreatedAt time.Time `db:"created_at"`
}
