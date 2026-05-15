package model

import "github.com/google/uuid"

type OrderItem struct {
	PartUUID uuid.UUID
	PartType PartType
	Price    int64
}
