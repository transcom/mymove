package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAddress creates a single Address and associated service member.
func MakePrimeContractor(db *pop.Connection, assertions Assertions) models.Contractor {
	contractor := models.Contractor{
		Name:           "Prime McPrimeContractor",
		Type:           "PRIME",
		ContractNumber: "HTC711-20-D-R030",
	}

	mergeModels(&contractor, assertions.Contractor)

	mustCreate(db, &contractor)

	return contractor
}
