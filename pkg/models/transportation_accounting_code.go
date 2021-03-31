package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// TransportationAccountingCode model struct that represents transportation accounting codes
type TransportationAccountingCode struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TAC       string    `json:"tac" db:"tac"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
