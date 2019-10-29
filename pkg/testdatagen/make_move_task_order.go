package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMoveTaskOrder creates a move task order
func MakeMoveTaskOrder(db *pop.Connection, assertions Assertions) models.MoveTaskOrder {
	id, _ := uuid.NewV4()
	mto := models.MoveTaskOrder{
		ID: id,
	}

	mergeModels(&mto, assertions.MoveTaskOrder)

	mustCreate(db, &mto)

	return mto
}

// MakeDefaultMoveTaskOrder makes an MoveTaskOrder with default values
func MakeDefaultMoveTaskOrder(db *pop.Connection) models.MoveTaskOrder {
	return MakeMoveTaskOrder(db, Assertions{})
}
