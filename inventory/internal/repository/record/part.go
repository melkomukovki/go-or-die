package record

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID `db:"uuid"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	PartType      string    `db:"part_type"`
	Price         int64     `db:"price"`
	StockQuantity int64     `db:"stock_quantity"`
	CreatedAt     time.Time `db:"created_at"`
}
