package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakePPM creates a single Personally Procured Move and its associated Move and Orders
func MakePPM(db *pop.Connection) (models.PersonallyProcuredMove, error) {
	move, err := MakeMove(db)
	if err != nil {
		return models.PersonallyProcuredMove{}, err
	}

	shirt := internalmessages.TShirtSizeM
	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)

	ppm, verrs, err := move.CreatePPM(db,
		&shirt,
		models.Int64Pointer(8000),
		models.StringPointer("estimate incentive"),
		models.TimePointer(DateInsidePeakRateCycle),
		models.StringPointer("72017"),
		models.BoolPointer(false),
		nil,
		models.StringPointer("60605"),
		models.BoolPointer(false),
		nil,
		models.StringPointer("estimate sit"),
		true,
		&advance,
	)

	if verrs.HasAny() || err != nil {
		return models.PersonallyProcuredMove{}, err
	}

	return *ppm, nil
}

// MakePPMData creates 3 PPMs (and in turn a move and set of orders for each)
func MakePPMData(db *pop.Connection) {
	for i := 0; i < 3; i++ {
		MakePPM(db)
	}
}
