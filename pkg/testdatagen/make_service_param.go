package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceParam creates a single ServiceParam
func MakeServiceParam(db *pop.Connection, assertions Assertions) models.ServiceParam {
	serviceParam := models.ServiceParam{}
	// Overwrite values with those from assertions
	mergeModels(&serviceParam, assertions.ServiceParam)

	mustCreate(db, &serviceParam)

	return serviceParam
}

// MakeDefaultServiceParam makes a ServiceParam with default values
func MakeDefaultServiceParam(db *pop.Connection) models.ServiceParam {
	return MakeServiceParam(db, Assertions{})
}
