package testdatagen

import (
	"database/sql"
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReContractYear creates a single ReContractYear with associations
func MakeReContractYear(db *pop.Connection, assertions Assertions) models.ReContractYear {
	reContract := assertions.ReContractYear.Contract
	if isZeroUUID(reContract.ID) {
		reContract = FetchOrMakeReContract(db, assertions)
	}

	reContractYear := models.ReContractYear{
		ContractID:           reContract.ID,
		Name:                 "Base Period Year 1",
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

func FetchOrMakeReContractYear(db *pop.Connection, assertions Assertions) models.ReContractYear {
	var existingContractYear models.ReContractYear
	if !assertions.ReContractYear.StartDate.IsZero() && !assertions.ReContractYear.EndDate.IsZero() {
		// converts the start and end dates into a daterange (inclusive bounds) and checks if
		// it overlaps with the daterange built from the given assertion dates
		err := db.Eager("Contract").Where(
			"daterange(start_date, end_date, '[]') && daterange(?, ?, '[]')",
			assertions.ReContractYear.StartDate, assertions.ReContractYear.EndDate,
		).First(&existingContractYear)
		if err != nil && err != sql.ErrNoRows {
			log.Panic("unexpected query error looking for overlapping ReContractYear", err)
		}
		if existingContractYear.ID != uuid.Nil {
			return existingContractYear
		}
	}
	return MakeReContractYear(db, assertions)
}
