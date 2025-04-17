package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReDomesticServiceArea creates a single ReDomesticServiceArea
func MakeReDomesticServiceArea(db *pop.Connection, assertions Assertions) models.ReDomesticServiceArea {
	var reContract models.ReContract

	if assertions.ReDomesticServiceArea.ContractID != uuid.Nil || assertions.ReDomesticServiceArea.Contract.ID != uuid.Nil || assertions.ReDomesticServiceArea.Contract.Code != "" {
		reContract = assertions.ReDomesticServiceArea.Contract
	} else if assertions.ReContract.ID != uuid.Nil || assertions.ReContract.Code != "" {
		reContract = assertions.ReContract
	}

	if isZeroUUID(reContract.ID) {
		assertions.ReContract = reContract
		reContract = FetchOrMakeReContract(db, assertions)
	}

	reDomesticServiceArea := models.ReDomesticServiceArea{
		ContractID:       reContract.ID,
		Contract:         reContract,
		ServiceArea:      "056",
		ServicesSchedule: 3,
		SITPDSchedule:    3,
	}

	// Overwrite values with those from assertions
	mergeModels(&reDomesticServiceArea, assertions.ReDomesticServiceArea)

	mustCreate(db, &reDomesticServiceArea, assertions.Stub)

	return reDomesticServiceArea
}

func FetchOrMakeReDomesticServiceArea(db *pop.Connection, assertions Assertions) models.ReDomesticServiceArea {

	var reDomesticServiceArea models.ReDomesticServiceArea
	if assertions.ReDomesticServiceArea.ServiceArea != "" {
		err := db.Where("re_domestic_service_areas.service_area = ?", assertions.ReDomesticServiceArea.ServiceArea).First(&reDomesticServiceArea)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		}
	} else {
		err := db.Where("re_domestic_service_areas.contract_id = ? AND re_domestic_service_areas.service_area = ?", assertions.ReDomesticServiceArea.ContractID, assertions.ReDomesticServiceArea.ServiceArea).First(&reDomesticServiceArea)
		if err != nil && err != sql.ErrNoRows {
			log.Panic(err)
		}
	}

	if assertions.ReDomesticServiceArea.ServiceArea == "" {
		return MakeReDomesticServiceArea(db, assertions)
	}

	if reDomesticServiceArea.ID == uuid.Nil {
		return MakeReDomesticServiceArea(db, assertions)
	}

	return reDomesticServiceArea
}

// MakeDefaultReDomesticServiceArea makes a single ReDomesticServiceArea with default values
func MakeDefaultReDomesticServiceArea(db *pop.Connection) models.ReDomesticServiceArea {
	return MakeReDomesticServiceArea(db, Assertions{})
}
