package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type ReportViolation struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	ReportID    uuid.UUID        `json:"report_id" db:"report_id"`
	Report      EvaluationReport `belongs_to:"evaluation_report" fk_id:"report_id"`
	Violation   PWSViolation     `belongs_to:"pws_violations" fk_id:"violation_id"`
	ViolationID uuid.UUID        `json:"violation_id" db:"violation_id"`
	GsrAppeals  GsrAppeals       `json:"gsr_appeal,omitempty" fk_id:"report_violation_id" has_many:"gsr_appeal"`
}

// TableName overrides the table name used by Pop.
func (r ReportViolation) TableName() string {
	return "report_violations"
}

type ReportViolations []ReportViolation

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReportViolation) Validate(_ *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator

	verrs := validate.Validate(vs...)

	return verrs, nil
}
