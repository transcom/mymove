package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"
)

type EvaluationReportType string

const (
	EvaluationReportTypeDataReview  EvaluationReportType = "DATA_REVIEW"
	EvaluationReportTypePhysical    EvaluationReportType = "PHYSICAL"
	EvaluationReportTypeDataVirtual EvaluationReportType = "VIRTUAL"
)

type EvaluationReportLocationType string

const (
	EvaluationReportLocationTypeOrigin      EvaluationReportLocationType = "ORIGIN"
	EvaluationReportLocationTypeDestination EvaluationReportLocationType = "DESTINATION"
	EvaluationReportLocationTypeOther       EvaluationReportLocationType = "OTHER"
)

type EvaluationReport struct {
	ID                      uuid.UUID                     `json:"id" db:"id"`
	Shipment                *MTOShipment                  `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	ShipmentID              *uuid.UUID                    `json:"shipment_id" db:"shipment_id"`
	InspectionDate          *time.Time                    `json:"inspection_date" db:"inspection_date"`
	Type                    *EvaluationReportType         `json:"type" db:"type"`
	TravelTimeMinutes       *int                          `json:"travel_time_minutes" db:"travel_time_minutes"`
	Location                *EvaluationReportLocationType `json:"location" db:"location"`
	LocationDescription     *string                       `json:"location_description" db:"location_description"`
	ObservedDate            *time.Time                    `json:"observed_date" db:"observed_date"`
	EvaluationLengthMinutes *int                          `json:"evaluation_length_minutes" db:"evaluation_length_minutes"`
	ViolationsObserved      *bool                         `json:"violations_observed" db:"violations_observed"`
	Remarks                 *string                       `json:"remarks" db:"remarks"`
	SubmittedAt             *time.Time                    `json:"submitted_at" db:"submitted_at"`
	CreatedAt               time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time                     `json:"updated_at" db:"updated_at"`
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
