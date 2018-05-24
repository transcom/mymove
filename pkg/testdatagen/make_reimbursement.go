package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeReimbursement creates a single Reimbursement
func MakeReimbursement(db *pop.Connection) (models.Reimbursement, error) {

	reimbursement, verrs, err := move.CreateReimbursement(db,
		&shirt,
		models.Int64Pointer(8000),
		models.StringPointer("estimate incentive"),
		models.TimePointer(DateInsidePeakRateCycle),
		models.StringPointer("72017"),
		models.BoolPointer(false),
		nil,
		models.StringPointer("60605"),
		models.BoolPointer(false),
		nil,
		true,
		&advance,
	)

	if verrs.HasAny() || err != nil {
		return models.Reimbursement{}, err
	}

	return *ppm, nil
}

// MakeReimbursementData creates 5 Reimbursements
func MakeReimbursementData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeReimbursement(db)
	}
}
