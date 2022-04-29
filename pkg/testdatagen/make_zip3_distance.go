package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeZip3Distance creates a single Zip3Distance
func MakeZip3Distance(db *pop.Connection, assertions Assertions) models.Zip3Distance {
	Zip3Distance := models.Zip3Distance{
		FromZip3:      "010",
		ToZip3:        "011",
		DistanceMiles: 24,
	}

	// Overwrite values with those from assertions
	mergeModels(&Zip3Distance, assertions.Zip3Distance)

	mustCreate(db, &Zip3Distance, assertions.Stub)

	return Zip3Distance
}

// MakeDefaultZip3Distance makes a single Zip3Distance with default values
func MakeDefaultZip3Distance(db *pop.Connection) models.Zip3Distance {
	return MakeZip3Distance(db, Assertions{})
}
