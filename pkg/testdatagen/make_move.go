package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMove creates a single Move and associated set of Orders
func MakeMove(db *pop.Connection, assertions Assertions) models.Move {

	// Create new Orders if not provided
	orders := assertions.Order
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.Order.ID) {
		orders = MakeOrder(db, assertions)
	}

	defaultMoveType := models.SelectedMoveTypePPM
	selectedMoveType := assertions.Move.SelectedMoveType
	if selectedMoveType == nil {
		selectedMoveType = &defaultMoveType
	}
	move := models.Move{
		Orders:           orders,
		OrdersID:         orders.ID,
		SelectedMoveType: selectedMoveType,
		Status:           models.MoveStatusDRAFT,
		Locator:          models.GenerateLocator(),
		Show:             setShow(assertions.Move.Show),
	}

	// Overwrite values with those from assertions
	mergeModels(&move, assertions.Move)

	mustCreate(db, &move)

	return move
}

// MakeMoveWithoutMoveType creates a single Move and associated set of Orders, but without a chosen move type
func MakeMoveWithoutMoveType(db *pop.Connection, assertions Assertions) models.Move {

	// Create new Orders if not provided
	orders := assertions.Order
	if isZeroUUID(assertions.Order.ID) {
		orders = MakeOrder(db, assertions)
	}

	move := models.Move{
		Orders:   orders,
		OrdersID: orders.ID,
		Status:   models.MoveStatusDRAFT,
		Locator:  models.GenerateLocator(),
		Show:     setShow(assertions.Move.Show),
	}

	// Overwrite values with those from assertions
	mergeModels(&move, assertions.Move)

	mustCreate(db, &move)

	return move
}

// MakeDefaultMove makes a Move with default values
func MakeDefaultMove(db *pop.Connection) models.Move {
	return MakeMove(db, Assertions{})
}

// MakeMoveData created 5 Moves (and in turn a set of Orders for each)
func MakeMoveData(db *pop.Connection) {
	for i := 0; i < 3; i++ {
		MakeDefaultMove(db)
	}

	for i := 0; i < 2; i++ {
		move := MakeDefaultMove(db)
		move.Approve()
		db.ValidateAndUpdate(&move)
	}
}

func setShow(assertionShow *bool) *bool {
	show := swag.Bool(true)
	if assertionShow != nil {
		show = assertionShow
	}
	return show
}
