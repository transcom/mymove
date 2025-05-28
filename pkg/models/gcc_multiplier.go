package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// GCCMultiplier represents the multipliers that apply to PPM incentives
type GCCMultiplier struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Multiplier float64   `json:"multiplier" db:"multiplier"`
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (g GCCMultiplier) TableName() string {
	return "gcc_multipliers"
}
