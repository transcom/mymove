package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"
)

type EvaluationReportType string

const (
	EvaluationReportTypeShipment   EvaluationReportType = "SHIPMENT"
	EvaluationReportTypeCounseling EvaluationReportType = "COUNSELING"
)

type EvaluationReportInspectionType string

const (
	EvaluationReportInspectionTypeDataReview  EvaluationReportInspectionType = "DATA_REVIEW"
	EvaluationReportInspectionTypePhysical    EvaluationReportInspectionType = "PHYSICAL"
	EvaluationReportInspectionTypeDataVirtual EvaluationReportInspectionType = "VIRTUAL"
)

type EvaluationReportLocationType string

const (
	EvaluationReportLocationTypeOrigin      EvaluationReportLocationType = "ORIGIN"
	EvaluationReportLocationTypeDestination EvaluationReportLocationType = "DESTINATION"
	EvaluationReportLocationTypeOther       EvaluationReportLocationType = "OTHER"
)

type EvaluationReport struct {
	ID                      uuid.UUID                       `json:"id" db:"id"`
	OfficeUser              OfficeUser                      `belongs_to:"office_users" fk_id:"office_user_id"`
	OfficeUserID            uuid.UUID                       `db:"office_user_id"`
	Move                    Move                            `belongs_to:"moves" fk_id:"move_id"`
	MoveID                  uuid.UUID                       `db:"move_id"`
	Shipment                *MTOShipment                    `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	ShipmentID              *uuid.UUID                      `json:"shipment_id" db:"shipment_id"`
	Type                    EvaluationReportType            `json:"type" db:"type"`
	InspectionDate          *time.Time                      `json:"inspection_date" db:"inspection_date"`
	InspectionType          *EvaluationReportInspectionType `json:"inspection_type" db:"inspection_type"`
	TravelTimeMinutes       *int                            `json:"travel_time_minutes" db:"travel_time_minutes"`
	Location                *EvaluationReportLocationType   `json:"location" db:"location"`
	LocationDescription     *string                         `json:"location_description" db:"location_description"`
	ObservedDate            *time.Time                      `json:"observed_date" db:"observed_date"`
	EvaluationLengthMinutes *int                            `json:"evaluation_length_minutes" db:"evaluation_length_minutes"`
	ViolationsObserved      *bool                           `json:"violations_observed" db:"violations_observed"`
	Remarks                 *string                         `json:"remarks" db:"remarks"`
	SubmittedAt             *time.Time                      `json:"submitted_at" db:"submitted_at"`
	DeletedAt               *time.Time                      `db:"deleted_at"`
	CreatedAt               time.Time                       `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time                       `json:"updated_at" db:"updated_at"`
}

// EvaluationReports is not required by pop and may be deleted
type EvaluationReports []EvaluationReport

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *EvaluationReport) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return nil, nil
}

func (r *EvaluationReport) TableName() string {
	return "evaluation_reports"
}
