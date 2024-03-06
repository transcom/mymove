package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakePPM creates a single Personally Procured Move and its associated Move and Orders
func MakePPM(db *pop.Connection, assertions Assertions) models.PersonallyProcuredMove {

	// Create new Move if not provided
	move := assertions.PersonallyProcuredMove.Move
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.PersonallyProcuredMove.MoveID) {
		move = makeMove(db, assertions)
	}

	ppm := models.PersonallyProcuredMove{
		Move:                          move,
		MoveID:                        move.ID,
		WeightEstimate:                models.PoundPointer(8000),
		OriginalMoveDate:              models.TimePointer(DateInsidePeakRateCycle),
		PickupPostalCode:              models.StringPointer("72017"),
		HasAdditionalPostalCode:       models.BoolPointer(false),
		AdditionalPickupPostalCode:    nil,
		DestinationPostalCode:         models.StringPointer("60605"),
		HasSit:                        models.BoolPointer(false),
		DaysInStorage:                 nil,
		Status:                        models.PPMStatusDRAFT,
		EstimatedStorageReimbursement: models.StringPointer("estimate sit"),
	}

	// Overwrite values with those from assertions
	mergeModels(&ppm, assertions.PersonallyProcuredMove)

	mustCreate(db, &ppm, assertions.Stub)
	return ppm
}

// MakeDefaultPPM makes a PPM with default values
func MakeDefaultPPM(db *pop.Connection) models.PersonallyProcuredMove {
	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
	MustSave(db, &advance)

	return MakePPM(db, Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			Advance:             &advance,
			AdvanceID:           &advance.ID,
			HasRequestedAdvance: true,
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
	})
}
