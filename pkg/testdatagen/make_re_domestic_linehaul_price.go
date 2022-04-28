package testdatagen

import (
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
