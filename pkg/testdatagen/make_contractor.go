package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeContractor creates a single Contractor.
func MakeContractor(db *pop.Connection, assertions Assertions) models.Contractor {

	var contractor models.Contractor

	if assertions.Contractor.Name == "" {
		assertions.Contractor.Name  = DefaultContractName
	}

	if assertions.Contractor.ContractNumber != ""  {
		assertions.Contractor.ContractNumber = DefaultContractCode
	}

	if assertions.Contractor.Type != "" {
		assertions.Contractor.Type = DefaultContractType
	}

	// Overwrite values with those from assertions
	mergeModels(&contractor, assertions.Document)

	mustCreate(db, &contractor)

	return contractor
}

// MakeDefaultContractor returns a Contractor with default values
func MakeDefaultContractor(db *pop.Connection) models.Contractor {
	return MakeContractor(db, Assertions{})
}

