package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReContract creates a single ReContract
func MakeReContract(db *pop.Connection, assertions Assertions) models.ReContract {
	reContract := models.ReContract{
		Code: DefaultContractCode,
		Name: "Test Contract",
	}

	// Overwrite values with those from assertions
	mergeModels(&reContract, assertions.ReContract)

	mustCreate(db, &reContract, assertions.Stub)

	return reContract
}

// MakeDefaultReContract makes a single ReContract with default values
func MakeDefaultReContract(db *pop.Connection) models.ReContract {
	return MakeReContract(db, Assertions{})
}
