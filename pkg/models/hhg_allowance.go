package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// HHGAllowance the allowance in weights for a given pay grade
type HHGAllowance struct {
	ID                            uuid.UUID `json:"id" db:"id"`
	PayGradeID                    uuid.UUID `json:"pay_grade_id" db:"pay_grade_id"`
	PayGrade                      PayGrade  `belongs_to:"pay_grades" fk_id:"pay_grade_id"`
	TotalWeightSelf               int       `json:"total_weight_self" db:"total_weight_self"`
	TotalWeightSelfPlusDependents int       `json:"total_weight_self_plus_dependents" db:"total_weight_self_plus_dependents"`
	ProGearWeight                 int       `json:"pro_gear_weight" db:"pro_gear_weight"`
	ProGearWeightSpouse           int       `json:"pro_gear_weight_spouse" db:"pro_gear_weight_spouse"`
	CreatedAt                     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at" db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" method
func (h HHGAllowance) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "PayGradeID", Field: h.PayGradeID},
		// Make sure we don't somehow get a negative value
		&validators.IntIsGreaterThan{Name: "TotalWeightSelf", Field: h.TotalWeightSelf, Compared: -1},
		&validators.IntIsGreaterThan{Name: "TotalWeightSelfPlusDependents", Field: h.TotalWeightSelfPlusDependents, Compared: -1},
		&validators.IntIsGreaterThan{Name: "ProGearWeight", Field: h.ProGearWeight, Compared: -1},
		&validators.IntIsGreaterThan{Name: "ProGearWeightSpouse", Field: h.ProGearWeightSpouse, Compared: -1},
	), nil
}

// HHGAllowances is a slice of HHGAllowance
type HHGAllowances []HHGAllowance

// TableName overrides the table name used by Pop.
func (h HHGAllowance) TableName() string {
	return "hhg_allowances"
}
