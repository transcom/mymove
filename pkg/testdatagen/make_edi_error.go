package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeEdiError creates a single EdiError
func MakeEdiError(db *pop.Connection, assertions Assertions) models.EdiError {
	ediError := models.EdiError{
		PaymentRequestID: assertions.EdiError.PaymentRequestID,
		Code:             assertions.EdiError.Code,
		Description:      assertions.EdiError.Description,
		EDIType:          assertions.EdiError.EDIType,
	}

	// Overwrite values with those from assertions
	mergeModels(&ediError, assertions.EdiError)

	mustCreate(db, &ediError, assertions.Stub)

	return ediError
}
