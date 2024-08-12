package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// AppealStatus represents the status of an appeal made by a GSR
type AppealStatus string

const (
	AppealStatusSustained string = "SUSTAINED"
	AppealStatusRejected  string = "REJECTED"
)

type GsrAppeal struct {
	ID                      uuid.UUID         `json:"id" db:"id"`
	EvaluationReportID      *uuid.UUID        `json:"evaluation_report_id" db:"evaluation_report_id"`
	EvaluationReport        *EvaluationReport `belongs_to:"evaluation_reports" fk_id:"evaluation_report_id"`
	ReportViolationID       *uuid.UUID        `json:"report_violation_id" db:"report_violation_id"`
	ReportViolation         *ReportViolation  `belongs_to:"report_violations" fk_id:"report_violations"`
	OfficeUserID            uuid.UUID         `json:"office_user_id" db:"office_user_id"`
	OfficeUser              OfficeUser        `belongs_to:"office_users" fk_id:"office_users"`
	IsSeriousIncidentAppeal *bool             `json:"is_serious_incident_appeal" db:"is_serious_incident_appeal"`
	AppealStatus            AppealStatus      `json:"appeal_status" db:"appeal_status"`
	Remarks                 string            `json:"remarks" db:"remarks"`
	CreatedAt               time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time        `json:"deleted_at" db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (mh GsrAppeal) TableName() string {
	return "gsr_appeals"
}
