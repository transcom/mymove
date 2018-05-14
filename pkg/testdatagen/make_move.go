package testdatagen

import (
	"github.com/gobuffalo/pop"

	"errors"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

// MakeMove creates a single Move and associated set of Orders
func MakeMove(db *pop.Connection) (models.Move, error) {
	orders, err := MakeOrder(db)
	if err != nil {
		return models.Move{}, err
	}

	var selectedType = internalmessages.SelectedMoveTypePPM
	move, ok := orders.CreateNewMove(db, zap.L(), &selectedType)
	if !ok {
		return models.Move{}, errors.New("failed to create move")
	}
	move.Orders = orders
	return *move, nil
}

// MakeMoveData created 5 Moves (and in turn a set of Orders for each)
func MakeMoveData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeMove(db)
	}
}
