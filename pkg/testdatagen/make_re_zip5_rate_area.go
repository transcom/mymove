package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReZip5RateArea creates a single ReZip5RateArea
func MakeReZip5RateArea(db *pop.Connection, assertions Assertions) models.ReZip5RateArea {
	reContract := assertions.ReContract
	if isZeroUUID(reContract.ID) {
		reContract = FetchOrMakeReContract(db, assertions)
	}

	rateArea := assertions.ReRateArea
	if isZeroUUID(rateArea.ID) {
		rateArea = FetchOrMakeReRateArea(db, assertions)
	}

	reZip5RateArea := models.ReZip5RateArea{
		Contract:   reContract,
		ContractID: reContract.ID,
		Zip5:       "32102",
		RateArea:   rateArea,
		RateAreaID: rateArea.ID,
	}

	// Overwrite values with those from assertions
	mergeModels(&reZip5RateArea, assertions.ReZip5RateArea)

	mustCreate(db, &reZip5RateArea, assertions.Stub)

	return reZip5RateArea
}
