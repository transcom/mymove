package testdatagen

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReService creates a single ReService
func MakeReService(db *pop.Connection, assertions Assertions) models.ReService {
	reService := models.ReService{
		Code: DefaultServiceCode,
		Name: "Test Service",
	}

	// Overwrite values with those from assertions
	mergeModels(&reService, assertions.ReService)

	mustCreate(db, &reService, assertions.Stub)

	return reService
}

// FetchOrMakeReService returns the ReService for a given service code, or creates one if
// the service code does not exist yet.
func FetchOrMakeReService(db *pop.Connection, assertions Assertions) models.ReService {
	var existingReServices models.ReServices
	code := DefaultServiceCode
	if assertions.ReService.Code != "" {
		code = string(assertions.ReService.Code)
	}
	err := db.Where("code = ?", code).All(&existingReServices)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}

	if len(existingReServices) == 0 {
		return MakeReService(db, assertions)
	}

	return existingReServices[0]
}

// MakeDefaultReService makes a single ReService with default values
func MakeDefaultReService(db *pop.Connection) models.ReService {
	return MakeReService(db, Assertions{})
}
