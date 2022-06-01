package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReContractYear creates a single ReContractYear with associations
func MakeReContractYear(db *pop.Connection, assertions Assertions) models.ReContractYear {
	reContract := assertions.ReContractYear.Contract
	if isZeroUUID(reContract.ID) {
		reContract = MakeReContract(db, assertions)
	}

	reContractYear := models.ReContractYear{
		ContractID:           reContract.ID,
		Name:                 "Test Contract Year",
		StartDate:            time.Date(TestYear, time.January, 1, 0, 0, 0, 0, time.UTC),
		EndDate:              time.Date(TestYear, time.December, 31, 0, 0, 0, 0, time.UTC),
		Escalation:           1.0197,
		EscalationCompounded: 1.04071,
		Contract:             reContract,
	}

	// Overwrite values with those from assertions
	mergeModels(&reContractYear, assertions.ReContractYear)

	mustCreate(db, &reContractYear, assertions.Stub)

	return reContractYear
}

// MakeDefaultReContractYear makes a single ReContractYear with default values
func MakeDefaultReContractYear(db *pop.Connection) models.ReContractYear {
	return MakeReContractYear(db, Assertions{})
}
