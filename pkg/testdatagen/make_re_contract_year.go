package testdatagen

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofrs/uuid"

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

func FetchOrMakeReContractYear(db *pop.Connection, assertions Assertions) models.ReContractYear {
	var existingContractYear models.ReContractYear
	if !assertions.ReContractYear.StartDate.IsZero() && !assertions.ReContractYear.EndDate.IsZero() {
		err := db.Where("start_date = ? AND end_date = ?", assertions.ReContractYear.StartDate, assertions.ReContractYear.EndDate).First(&existingContractYear)
		if err != nil && err != sql.ErrNoRows {
			log.Panic("unexpected query error looking for existing ReContractYear by start and end dates", err)
		}

		if existingContractYear.ID != uuid.Nil {
			return existingContractYear
		}
	}

	return MakeReContractYear(db, assertions)
}
