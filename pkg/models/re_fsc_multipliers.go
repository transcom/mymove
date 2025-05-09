package models

import (
	"time"

	"github.com/gofrs/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
)

// Note: Multiplier is a pointer to avoid copying a struct that contains a sync.Mutex.
type FscMultiplier struct {
	ID         uuid.UUID        `json:"id" db:"id" rw:"r"`
	LowWeight  int              `json:"low_weight" rw:"r"`
	HighWeight int              `json:"high_weight" rw:"r"`
	Multiplier *decimal.Decimal `json:"multiplier"  rw:"r"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt  time.Time        `json:"updated_at" db:"updated_at" rw:"r"`
}

func (f FscMultiplier) TableName() string {
	return "re_fsc_multipliers"
}
