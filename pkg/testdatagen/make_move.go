package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeMove creates a single Move and associated set of Orders
func MakeMove(db *pop.Connection) (models.Move, error) {
	orders, err := MakeOrder(db)
	if err != nil {
		return models.Move{}, err
	}

	var selectedType = internalmessages.SelectedMoveTypePPM

	move, verrs, err := orders.CreateNewMove(db, &selectedType)
	if verrs.HasAny() || err != nil {
		return models.Move{}, err
	}
	// Can we move this into CreateNewMove
	// (pivotal: https://www.pivotaltracker.com/story/show/157610355)
	move.Orders = orders

	return *move, nil
}

// MakeMoveData created 5 Moves (and in turn a set of Orders for each)
func MakeMoveData(db *pop.Connection) {
	for i := 0; i < 3; i++ {
		MakeMove(db)
	}

	for i := 0; i < 2; i++ {
		move, _ := MakeMove(db)
		move.Approve()
		db.ValidateAndUpdate(&move)
	}
}
