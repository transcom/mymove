package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// MoveTaskOrder is an object representing a move task order
type MoveTaskOrder struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	MoveID       uuid.UUID   `json:"move_id" db:"move_id"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	ActualWeight *unit.Pound `json:"actual_weight" db:"actual_weight"`
}

// FetchMoveTaskOrder fetches the move task order by id
// We might want to consider moving this to its own service object
func FetchMoveTaskOrder(db *pop.Connection, id uuid.UUID) (*MoveTaskOrder, error) {
	mto := MoveTaskOrder{}
	err := db.Where("id = ?", id).First(&mto)
	if err != nil {
		return &mto, err
	}
	return &mto, nil
}
