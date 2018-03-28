package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeMove creates a single Move and associated User.
func MakeMove(db *pop.Connection) (models.Move, error) {
	var move models.Move

	user, err := MakeUser(db)
	if err != nil {
		return move, err
	}
	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move = models.Move{
		UserID:           user.ID,
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

// MakeMoveData created 5 Moves (and in turn a User for each)
func MakeMoveData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeMove(db)
	}
}
