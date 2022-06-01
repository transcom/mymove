package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// SetupServiceAreaRateArea sets up contract, service area, rate area, zip3
// returns contractYear, serviceArea, rateArea, reZip3
func SetupServiceAreaRateArea(db *pop.Connection, assertions Assertions) (models.ReContractYear, models.ReDomesticServiceArea, models.ReRateArea, models.ReZip3) {
	contractYear := models.ReContractYear{
		Escalation:           1.0185,
		EscalationCompounded: 1.04082,
		StartDate:            time.Date(GHCTestYear, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:              time.Date(GHCTestYear, time.December, 31, 0, 0, 0, 0, time.UTC),
	}

	mergeModels(&contractYear, assertions.ReContractYear)

	contractYear = MakeReContractYear(db, Assertions{ReContractYear: contractYear})

	serviceArea := models.ReDomesticServiceArea{
		Contract:    contractYear.Contract,
		ServiceArea: "042",
	}

	mergeModels(&serviceArea, assertions.ReDomesticServiceArea)

	serviceArea = MakeReDomesticServiceArea(db,
		Assertions{
			ReDomesticServiceArea: serviceArea,
		})

	rateArea := models.ReRateArea{
		ContractID: contractYear.Contract.ID,
		IsOconus:   false,
		Code:       "US47",
		Name:       "CA",
		Contract:   contractYear.Contract,
	}

	mergeModels(&rateArea, assertions.ReRateArea)

	rateArea = FetchOrMakeReRateArea(db, Assertions{
		ReContractYear: contractYear,
		ReRateArea:     rateArea,
	})

	reZip3 := models.ReZip3{
		ContractID:            contractYear.Contract.ID,
		Zip3:                  "940",
		BasePointCity:         "San Francisco",
		State:                 "CA",
		DomesticServiceAreaID: serviceArea.ID,
		RateAreaID:            &rateArea.ID,
		HasMultipleRateAreas:  false,
		Contract:              contractYear.Contract,
		DomesticServiceArea:   serviceArea,
		RateArea:              &rateArea,
	}

	mergeModels(&reZip3, assertions.ReZip3)

	reZip3 = MakeReZip3(db, Assertions{
		ReZip3: reZip3,
	})

	return contractYear, serviceArea, rateArea, reZip3
}
