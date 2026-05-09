package model

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

type PartFilter struct {
	UUIDs    []uuid.UUID
	PartType PartType
}
