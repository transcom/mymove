package testdatagen

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReDomesticServiceArea creates a single ReDomesticServiceArea
func MakeReDomesticServiceArea(db *pop.Connection, assertions Assertions) models.ReDomesticServiceArea {
	reContract := assertions.ReDomesticServiceArea.Contract
	if isZeroUUID(reContract.ID) {
		reContract = MakeReContract(db, assertions)
	}

	reDomesticServiceArea := models.ReDomesticServiceArea{
		ContractID:       reContract.ID,
		Contract:         reContract,
		ServiceArea:      "004",
		ServicesSchedule: 2,
		SITPDSchedule:    2,
	}

	// Overwrite values with those from assertions
	mergeModels(&reDomesticServiceArea, assertions.ReDomesticServiceArea)

	mustCreate(db, &reDomesticServiceArea, assertions.Stub)

	return reDomesticServiceArea
}

// MakeDefaultReDomesticServiceArea makes a single ReDomesticServiceArea with default values
func MakeDefaultReDomesticServiceArea(db *pop.Connection) models.ReDomesticServiceArea {
	return MakeReDomesticServiceArea(db, Assertions{})
}
