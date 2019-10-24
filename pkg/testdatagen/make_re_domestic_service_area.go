package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeReDomesticServiceArea creates a single ReDomesticServiceArea
func MakeReDomesticServiceArea(db *pop.Connection, assertions Assertions) models.ReDomesticServiceArea {
	reDomesticServiceArea := models.ReDomesticServiceArea{
		BasePointCity:    "Birmingham",
		State:            "AL",
		ServiceArea:      "004",
		ServicesSchedule: 2,
		SITPDSchedule:    2,
	}

	// Overwrite values with those from assertions
	mergeModels(&reDomesticServiceArea, assertions.ReDomesticServiceArea)

	mustCreate(db, &reDomesticServiceArea)

	return reDomesticServiceArea
}

// MakeDefaultReDomesticServiceArea makes a single ReDomesticServiceArea with default values
func MakeDefaultReDomesticServiceArea(db *pop.Connection) models.ReDomesticServiceArea {
	return MakeReDomesticServiceArea(db, Assertions{})
}
