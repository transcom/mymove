package testdatagen

import (
	"log"

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
	move := models.Move{
		OrdersID:         orders.ID,
		Orders:           orders,
		SelectedMoveType: &selectedType,
	}

	verrs, err := db.ValidateAndSave(&move)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return move, err
}

// MakeMoveData created 5 Moves (and in turn a set of Orders for each)
func MakeMoveData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeMove(db)
	}
}
