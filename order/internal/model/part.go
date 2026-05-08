package model

import "github.com/google/uuid"

type Part struct {
	UUID          uuid.UUID
	Name          string
	PartType      PartType
	Price         int64
	StockQuantity int64
}

type PartType string

const (
	PartTypeHull   PartType = "HULL"
	PartTypeEngine PartType = "ENGINE"
	PartTypeShield PartType = "SHIELD"
	PartTypeWeapon PartType = "WEAPON"
)
