package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// FetchOrMakeReRateArea returns the ReRateArea for a given rate area code, or creates one if
// the rate area does not exist yet.
func FetchOrMakeReRateArea(db *pop.Connection, assertions Assertions) models.ReRateArea {
	var existingReRateAreas models.ReRateAreas
	code := "US42"
	if assertions.ReRateArea.Code != "" {
		code = assertions.ReRateArea.Code
	}

	err := db.Where("code = ?", code).All(&existingReRateAreas)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if len(existingReRateAreas) == 0 {
		var contractYear models.ReContractYear
		if assertions.ReContractYear.ID == uuid.Nil {
			contractYear = MakeReContractYear(db, assertions)
		} else {
			contractYear = assertions.ReContractYear
		}

		rateArea := models.ReRateArea{
			ContractID: contractYear.Contract.ID,
			IsOconus:   false,
			Code:       code,
			Name:       "CA",
			Contract:   contractYear.Contract,
		}

		mergeModels(&rateArea, assertions.ReRateArea)

		MustSave(db, &rateArea)
		existingReRateAreas = append(existingReRateAreas, rateArea)
	}

	return existingReRateAreas[0]
}
