package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"

	"github.com/transcom/mymove/pkg/unit"
)

type PPMCloseout struct {
	PlannedMoveDate            *time.Time
	ActualMoveDate             *time.Time
	Miles                      *unit.Miles
	EstimatedWeight            *unit.Pound
	ActualWeight               *unit.Pound
	ProGearWeight              *unit.Pound
	GrossIncentive             *unit.Cents
	GCC                        *unit.Cents
	AOA                        *unit.Cents
	RemainingReimbursementOwed *unit.Cents
	HaulPrice                  *unit.Cents
	HaulFSC                    *unit.Cents
	DOP                        *unit.Cents
	DDP                        *unit.Cents
	PackUnpackPrice            *unit.Cents
	SITReimbursement           *unit.Cents
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (p PPMCloseout) Validate(_ *pop.Connection) (*validate.Errors, error) {
	milesInt := p.Miles.Int()
	return validate.Validate(
		&validators.TimeIsPresent{Name: "PlannedMoveDate", Field: *p.PlannedMoveDate},
		&validators.TimeIsPresent{Name: "ActualMoveDate", Field: *p.ActualMoveDate},
		&OptionalIntIsPositive{Name: "Miles", Field: &milesInt},
		&OptionalPoundIsNonNegative{Name: "EstimatedWeight", Field: p.EstimatedWeight},
		&OptionalPoundIsNonNegative{Name: "ActualWeight", Field: p.ActualWeight},
		&OptionalPoundIsNonNegative{Name: "ProGearWeight", Field: p.ProGearWeight},
		&OptionalCentIsPositive{Name: "GrossIncentive", Field: p.GrossIncentive},
		&OptionalCentIsPositive{Name: "GCC", Field: p.GCC},
		&OptionalCentIsPositive{Name: "AOA", Field: p.AOA},
		&OptionalCentIsPositive{Name: "RemainingReimbursementOwed", Field: p.RemainingReimbursementOwed},
		&OptionalCentIsPositive{Name: "HaulPrice", Field: p.HaulPrice},
		&OptionalCentIsPositive{Name: "HaulFSC", Field: p.HaulFSC},
		&OptionalCentIsPositive{Name: "DOP", Field: p.DOP},
		&OptionalCentIsPositive{Name: "DDP", Field: p.DDP},
		&OptionalCentIsPositive{Name: "PackUnpackPrice", Field: p.PackUnpackPrice},
		&OptionalCentIsPositive{Name: "SITReimbursement", Field: p.SITReimbursement},
	), nil
}