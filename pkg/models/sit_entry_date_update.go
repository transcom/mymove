package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type SITEntryDateUpdate struct {
	ID           uuid.UUID  `db:"id"`
	SITEntryDate *time.Time `db:"sit_entry_date"`
}
