package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func makeMinimalMovingExpense(db *pop.Connection, assertions Assertions) models.MovingExpense {
	ppmShipment := checkOrCreatePPMShipment(db, assertions)

	document := GetOrCreateDocument(db, assertions.MovingExpense.Document, assertions)

	newMovingExpense := models.MovingExpense{
		PPMShipmentID: ppmShipment.ID,
		PPMShipment:   ppmShipment,
		DocumentID:    document.ID,
		Document:      document,
	}

	// Overwrites model with data from assertions
	mergeModels(&newMovingExpense, assertions.MovingExpense)

	mustCreate(db, &newMovingExpense, assertions.Stub)

	return newMovingExpense
}

func makeMovingExpense(db *pop.Connection, assertions Assertions) models.MovingExpense {
	document := GetOrCreateDocumentWithUploads(db, assertions.MovingExpense.Document, assertions)
	packingMaterialType := models.MovingExpenseReceiptTypePackingMaterials
	amountPaid := unit.Cents(2345)

	fullAssertions := Assertions{
		MovingExpense: models.MovingExpense{
			DocumentID:        document.ID,
			Document:          document,
			MovingExpenseType: &packingMaterialType,
			Description:       models.StringPointer("Packing Peanuts"),
			PaidWithGTCC:      models.BoolPointer(true),
			Amount:            &amountPaid,
			MissingReceipt:    models.BoolPointer(false),
		},
	}

	mergeModels(&fullAssertions, assertions)

	return makeMinimalMovingExpense(db, fullAssertions)
}
