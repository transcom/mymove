package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReZip3 creates a single ReZip3
func MakeReZip3(db *pop.Connection, assertions Assertions) models.ReZip3 {
	reContract := assertions.ReZip3.Contract
	if isZeroUUID(reContract.ID) {
		reContract = FetchOrMakeReContract(db, assertions)
	}

	domesticServiceArea := assertions.ReZip3.DomesticServiceArea
	if isZeroUUID(domesticServiceArea.ID) {
		domesticServiceArea = MakeReDomesticServiceArea(db, assertions)
	}

	reZip3 := models.ReZip3{
		DomesticServiceAreaID: domesticServiceArea.ID,
		Contract:              reContract,
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

func FetchOrMakeReZip3(db *pop.Connection, assertions Assertions) models.ReZip3 {
	var contractID uuid.UUID
	if assertions.ReZip3.ContractID != uuid.Nil {
		contractID = assertions.ReZip3.ContractID
	} else if assertions.ReContract.ID != uuid.Nil {
		contractID = assertions.ReContract.ID
	}

	if contractID == uuid.Nil || assertions.ReZip3.Zip3 == "" {
		return MakeReZip3(db, assertions)
	}

	var reZip3 models.ReZip3
	err := db.Eager("Contract").Where("re_zip3s.contract_id = ? AND re_zip3s.zip3 = ?", contractID, assertions.ReZip3.Zip3).First(&reZip3)

	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if reZip3.ID == uuid.Nil {
		return MakeReZip3(db, assertions)
	}

	return reZip3
}

// MakeDefaultReZip3 makes a single ReZip3 with default values
func MakeDefaultReZip3(db *pop.Connection) models.ReZip3 {
	return MakeReZip3(db, Assertions{})
}
