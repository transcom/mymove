package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeDraftReimbursement creates a single draft status Reimbursement
func MakeDraftReimbursement(db *pop.Connection) (models.Reimbursement, error) {

	reimbursement := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)

	if verrs, err := db.ValidateAndSave(&reimbursement); verrs.HasAny() || err != nil {
		panic(fmt.Sprintf("validation erros: %v, saving error: %s", verrs, err))
	}

	return reimbursement, nil
}

// MakeRequestedReimbursement creates a single requested status Reimbursement
func MakeRequestedReimbursement(db *pop.Connection) (models.Reimbursement, error) {

	reimbursement := models.BuildRequestedReimbursement(2000, models.MethodOfReceiptGTCC)

	if verrs, err := db.ValidateAndSave(&reimbursement); verrs.HasAny() || err != nil {
		panic(fmt.Sprintf("validation erros: %v, saving error: %s", verrs, err))
	}

	return reimbursement, nil
}

// MakeReimbursementData creates 3 draft Reimbursements and 2 requested Reimbursements
func MakeReimbursementData(db *pop.Connection) {
	for i := 0; i < 3; i++ {
		MakeDraftReimbursement(db)
	}
	for i := 0; i < 2; i++ {
		MakeRequestedReimbursement(db)
	}
}
