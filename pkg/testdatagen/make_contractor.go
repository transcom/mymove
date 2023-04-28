package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// FetchOrMakeContractor creates a single Contractor.
func FetchOrMakeContractor(db *pop.Connection, assertions Assertions) models.Contractor {

	contractor := models.Contractor{
		Name:           DefaultContractName,
		ContractNumber: DefaultContractNumber,
		Type:           DefaultContractType,
	}

	if assertions.Contractor.ContractNumber != "" {
		contractor.ContractNumber = assertions.Contractor.ContractNumber
	}
	if assertions.Contractor.Name != "" {
		contractor.Name = assertions.Contractor.Name
	}
	if assertions.Contractor.Type != "" {
		contractor.Type = assertions.Contractor.Type
	}

	err := db.Q().Where(`contract_number=$1`, contractor.ContractNumber).First(&contractor)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return contractor
	}

	// Don't create multiple contractors of the same type
	err = db.Q().Where(`type=$1`, contractor.Type).First(&contractor)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return contractor
	}

	// Overwrite values with those from assertions
	mergeModels(&contractor, assertions.Contractor)

	mustCreate(db, &contractor, assertions.Stub)

	return contractor
}
