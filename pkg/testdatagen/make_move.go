package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
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
		contractor := FetchOrMakeContractor(db, assertions)
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

	return move
}

func setShow(assertionShow *bool) *bool {
	show := swag.Bool(true)
	if assertionShow != nil {
		show = assertionShow
	}
	return show
}
