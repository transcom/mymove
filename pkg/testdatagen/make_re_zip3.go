package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReZip3 creates a single ReZip3
func MakeReZip3(db *pop.Connection, assertions Assertions) models.ReZip3 {
	reContract := assertions.ReZip3.Contract
	if isZeroUUID(reContract.ID) {
		reContract = MakeReContract(db, assertions)
	}

	domesticServiceArea := assertions.ReZip3.DomesticServiceArea
	if isZeroUUID(domesticServiceArea.ID) {
		domesticServiceArea = MakeReDomesticServiceArea(db, assertions)
	}

	reZip3 := models.ReZip3{
		DomesticServiceAreaID: domesticServiceArea.ID,
		ContractID:            reContract.ID,
		Zip3:                  "350",
		BasePointCity:         "Memphis",
		State:                 "TN",
	}

	// Overwrite values with those from assertions
	mergeModels(&reZip3, assertions.ReZip3)

	mustCreate(db, &reZip3, assertions.Stub)

	return reZip3
}

// MakeDefaultReZip3 makes a single ReZip3 with default values
func MakeDefaultReZip3(db *pop.Connection) models.ReZip3 {
	return MakeReZip3(db, Assertions{})
}
