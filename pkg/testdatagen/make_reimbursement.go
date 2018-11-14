package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDraftReimbursement creates a single draft status Reimbursement
func MakeDraftReimbursement(db *pop.Connection) (models.Reimbursement, error) {

	reimbursement := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)

	mustCreate(db, &reimbursement)

	return reimbursement, nil
}

// MakeRequestedReimbursement creates a single requested status Reimbursement
func MakeRequestedReimbursement(db *pop.Connection) (models.Reimbursement, error) {

	reimbursement := models.BuildRequestedReimbursement(2000, models.MethodOfReceiptGTCC)

	mustCreate(db, &reimbursement)

	return reimbursement, nil
}
