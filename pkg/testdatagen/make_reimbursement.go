package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeRequestedReimbursement creates a single requested status Reimbursement
func MakeRequestedReimbursement(db *pop.Connection, assertions Assertions) models.Reimbursement {

	reimbursement := models.BuildRequestedReimbursement(2000, models.MethodOfReceiptGTCC)

	mustCreate(db, &reimbursement, assertions.Stub)

	return reimbursement
}

// MakeDefaultRequestedReimbursement makes a user with default values
func MakeDefaultRequestedReimbursement(db *pop.Connection) models.Reimbursement {
	return MakeRequestedReimbursement(db, Assertions{})
}
