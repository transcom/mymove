package testdatagen

import (
	"github.com/gobuffalo/pop"

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

	mustCreate(db, &reService)

	return reService
}

// MakeDefaultReService makes a single ReService with default values
func MakeDefaultReService(db *pop.Connection) models.ReService {
	return MakeReService(db, Assertions{})
}
