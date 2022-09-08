package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReDomesticLinehaulPrice creates a single ReDomesticLinehaulPrice
func MakeReDomesticLinehaulPrice(db *pop.Connection, assertions Assertions) models.ReDomesticLinehaulPrice {

	serviceArea := assertions.ReDomesticLinehaulPrice.DomesticServiceArea
	serviceAreaID := assertions.ReDomesticLinehaulPrice.DomesticServiceAreaID
	if isZeroUUID(serviceAreaID) {
		if isZeroUUID(serviceArea.ID) {
			serviceArea = MakeReDomesticServiceArea(db, assertions)
			assertions.ReDomesticLinehaulPrice.DomesticServiceAreaID = serviceArea.ID
		} else {
			assertions.ReDomesticLinehaulPrice.DomesticServiceAreaID = serviceArea.ID
		}
	}

	reContract := assertions.ReDomesticLinehaulPrice.Contract
	reContractID := assertions.ReDomesticLinehaulPrice.ContractID
	if isZeroUUID(reContractID) {
		if isZeroUUID(reContract.ID) {
			reContract = MakeReContract(db, assertions)
			assertions.ReDomesticLinehaulPrice.ContractID = reContract.ID
		} else {
			assertions.ReDomesticLinehaulPrice.ContractID = reContract.ID
		}
	}

	id, _ := uuid.NewV4()
	reDomesticLinehaulPrice := models.ReDomesticLinehaulPrice{
		ID:                    id,
		ContractID:            reContract.ID,
		WeightLower:           500,
		WeightUpper:           4999,
		MilesLower:            1001,
		MilesUpper:            1500,
		IsPeakPeriod:          false,
		DomesticServiceAreaID: serviceArea.ID,
		PriceMillicents:       5000, // $0.050
	}

	// Overwrite values with those from assertions
	mergeModels(&reDomesticLinehaulPrice, assertions.ReDomesticLinehaulPrice)

	mustCreate(db, &reDomesticLinehaulPrice, assertions.Stub)

	return reDomesticLinehaulPrice
}

// MakeDefaultReDomesticLinehaulPrice makes a single ReDomesticLinehaulPrice with default values
func MakeDefaultReDomesticLinehaulPrice(db *pop.Connection) models.ReDomesticLinehaulPrice {
	return MakeReDomesticLinehaulPrice(db, Assertions{})
}

func FetchOrMakeReDomesticLinehaulPrice(db *pop.Connection, assertions Assertions) models.ReDomesticLinehaulPrice {
	reDomesticLinehaulPrice := assertions.ReDomesticLinehaulPrice
	var existingReDomesticLinehaulPrice models.ReDomesticLinehaulPrice
	if reDomesticLinehaulPrice.ContractID != uuid.Nil &&
		reDomesticLinehaulPrice.DomesticServiceAreaID != uuid.Nil &&
		reDomesticLinehaulPrice.MilesLower > 0 &&
		reDomesticLinehaulPrice.MilesUpper > 0 &&
		reDomesticLinehaulPrice.WeightLower.Int() > 0 &&
		reDomesticLinehaulPrice.WeightUpper.Int() > 0 {

		err := db.Where("contract_id = ? AND domestic_service_area_id = ? AND miles_lower = ? AND miles_upper = ? AND weight_lower = ? and weight_upper = ? and is_peak_period = ?",
			reDomesticLinehaulPrice.ContractID, reDomesticLinehaulPrice.DomesticServiceAreaID, reDomesticLinehaulPrice.MilesLower, reDomesticLinehaulPrice.MilesUpper, reDomesticLinehaulPrice.WeightLower, reDomesticLinehaulPrice.WeightUpper, reDomesticLinehaulPrice.IsPeakPeriod).First(&existingReDomesticLinehaulPrice)
		if err != nil && err != sql.ErrNoRows {
			log.Panic("unexpected query error looking for existing ReDomesticLinehaulPrice", err)
		}

		if existingReDomesticLinehaulPrice.ID != uuid.Nil {
			return reDomesticLinehaulPrice
		}
	}

	return MakeReDomesticLinehaulPrice(db, assertions)
}
