package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// AppealStatus represents the status of an appeal made by a GSR
type AppealStatus string

// String is a string representation of a GSR Appeal Status
func (g AppealStatus) String() string {
	return string(g)
}

const (
	AppealStatusSustained AppealStatus = "SUSTAINED"
	AppealStatusRejected  AppealStatus = "REJECTED"
)

var validAppealStatus = []string{
	string(AppealStatusSustained),
	string(AppealStatusRejected),
}

type GsrAppeal struct {
	ID                      uuid.UUID         `json:"id" db:"id"`
	EvaluationReportID      uuid.UUID         `json:"evaluation_report_id" db:"evaluation_report_id"`
	EvaluationReport        *EvaluationReport `belongs_to:"evaluation_reports" fk_id:"evaluation_report_id"`
	ReportViolationID       *uuid.UUID        `json:"report_violation_id" db:"report_violation_id"`
	ReportViolation         *ReportViolation  `belongs_to:"report_violations" fk_id:"report_violations"`
	OfficeUserID            uuid.UUID         `json:"office_user_id" db:"office_user_id"`
	OfficeUser              *OfficeUser       `belongs_to:"office_users" fk_id:"office_users"`
	IsSeriousIncidentAppeal *bool             `json:"is_serious_incident_appeal" db:"is_serious_incident_appeal"`
	AppealStatus            AppealStatus      `json:"appeal_status" db:"appeal_status"`
	Remarks                 string            `json:"remarks" db:"remarks"`
	CreatedAt               time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time        `json:"deleted_at" db:"deleted_at"`
}

type GsrAppeals []GsrAppeal

// TableName overrides the table name used by Pop.
func (g GsrAppeal) TableName() string {
	return "gsr_appeals"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (g *GsrAppeal) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: g.OfficeUserID, Name: "OfficeUserID"},
		&validators.StringInclusion{Field: g.AppealStatus.String(), Name: "AppealStatus", List: validAppealStatus},
		&validators.StringIsPresent{Field: g.Remarks, Name: "Remarks"},
	), nil
}
