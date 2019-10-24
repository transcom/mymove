package testdatagen

import (
	"github.com/gofrs/uuid"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDefaultMoveTaskOrder creates a single MoveTaskOrder
func MakeDefaultMoveTaskOrder(db *pop.Connection) models.MoveTaskOrder {
	move := MakeDefaultMove(db)
	mto := models.MoveTaskOrder{
		ID:     uuid.Must(uuid.FromString("cccc8162-cb62-484c-87bb-05a27f41ad61")),
		MoveID: move.ID,
	}

	mustCreate(db, &mto)

	return mto
}
