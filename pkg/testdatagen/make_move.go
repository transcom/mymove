package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// makeMove creates a single Move and associated set of Orders
//
// Deprecated: use factory.BuildMove for all new code
func makeMove(db *pop.Connection, assertions Assertions) (models.Move, error) {

	// Create new Orders if not provided
	orders := assertions.Order
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.Order.ID) {
		var err error
		orders, err = makeOrder(db, assertions)
		if err != nil {
			return models.Move{}, err
		}
	}

	assertedReferenceID := assertions.Move.ReferenceID
	var referenceID string
	if assertedReferenceID == nil || *assertedReferenceID == "" {
		referenceID, _ = models.GenerateReferenceID(db)
	}

	var contractorID uuid.UUID
	moveContractorID := assertions.Move.ContractorID
	if moveContractorID == nil {
		contractor := fetchOrMakeContractor(db, assertions)
		contractorID = contractor.ID
	}

	ppmType := assertions.Move.PPMType
	if assertions.Move.PPMType == nil {
		partialType := "PARTIAL"
		ppmType = &partialType
	}

	move := models.Move{
		Orders:       orders,
		OrdersID:     orders.ID,
		PPMType:      ppmType,
		Status:       models.MoveStatusDRAFT,
		Locator:      models.GenerateLocator(),
		Show:         setShow(assertions.Move.Show),
		ContractorID: &contractorID,
		ReferenceID:  &referenceID,
	}

	// Overwrite values with those from assertions
	mergeModels(&move, assertions.Move)

	mustCreate(db, &move, assertions.Stub)

	return move, nil
}

func setShow(assertionShow *bool) *bool {
	show := models.BoolPointer(true)
	if assertionShow != nil {
		show = assertionShow
	}
	return show
}
