package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

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

	assertedReferenceID := assertions.Move.ReferenceID
	var referenceID string
	if assertedReferenceID == nil || *assertedReferenceID == "" {
		referenceID, _ = models.GenerateReferenceID(db)
	}

	var contractorID uuid.UUID
	moveContractorID := assertions.Move.ContractorID
	if moveContractorID == nil {
		contractor := MakeContractor(db, assertions)
		contractorID = contractor.ID
	}

	defaultMoveType := models.SelectedMoveTypePPM
	selectedMoveType := assertions.Move.SelectedMoveType
	if selectedMoveType == nil {
		selectedMoveType = &defaultMoveType
	}

	ppmType := assertions.Move.PPMType
	if assertions.Move.PPMType == nil {
		partialType := "PARTIAL"
		ppmType = &partialType
	}

	move := models.Move{
		Orders:           orders,
		OrdersID:         orders.ID,
		SelectedMoveType: selectedMoveType,
		PPMType:          ppmType,
		Status:           models.MoveStatusDRAFT,
		Locator:          models.GenerateLocator(),
		Show:             setShow(assertions.Move.Show),
		ContractorID:     &contractorID,
		ReferenceID:      &referenceID,
	}

	// Overwrite values with those from assertions
	mergeModels(&move, assertions.Move)

	mustCreate(db, &move, assertions.Stub)

	return move
}

// MakeMoveWithoutMoveType creates a single Move and associated set of Orders, but without a chosen move type
func MakeMoveWithoutMoveType(db *pop.Connection, assertions Assertions) models.Move {

	// Create new Orders if not provided
	orders := assertions.Order
	if isZeroUUID(assertions.Order.ID) {
		orders = MakeOrder(db, assertions)
	}

	var referenceID string
	assertedReferenceID := assertions.Move.ReferenceID
	if assertedReferenceID == nil || *assertedReferenceID == "" {
		referenceID, _ = models.GenerateReferenceID(db)
	}

	var contractorID uuid.UUID
	moveContractorID := assertions.Move.ContractorID
	if moveContractorID == nil {
		contractor := MakeContractor(db, assertions)
		contractorID = contractor.ID
	}

	move := models.Move{
		Orders:       orders,
		OrdersID:     orders.ID,
		Status:       models.MoveStatusDRAFT,
		Locator:      models.GenerateLocator(),
		Show:         setShow(assertions.Move.Show),
		ContractorID: &contractorID,
		ReferenceID:  &referenceID,
	}

	// Overwrite values with those from assertions
	mergeModels(&move, assertions.Move)

	mustCreate(db, &move, assertions.Stub)

	return move
}

// MakeAvailableMove makes a Move that is available to the prime at
// the time of its creation
func MakeAvailableMove(db *pop.Connection) models.Move {
	now := time.Now()
	move := MakeMove(db, Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
			Status:             models.MoveStatusAPPROVED,
		},
	})
	return move
}

// MakeDefaultMove makes a Move with default values
func MakeDefaultMove(db *pop.Connection) models.Move {
	return MakeMove(db, Assertions{})
}

func setShow(assertionShow *bool) *bool {
	show := swag.Bool(true)
	if assertionShow != nil {
		show = assertionShow
	}
	return show
}
