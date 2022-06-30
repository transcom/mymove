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
	ShipmentID              *uuid.UUID                    `db:"shipment_id"`
	InspectionDate          *time.Time                    `db:"inspection_date"`
	Type                    *EvaluationReportType         `db:"type"`
	TravelTimeMinutes       *int                          `db:"travel_time_minutes"`
	Location                *EvaluationReportLocationType `db:"location"`
	LocationDescription     *string                       `db:"location_description"`
	ObservedDate            *time.Time                    `db:"observed_date"`
	EvaluationLengthMinutes *int                          `db:"evaluation_length_minutes"`
	ViolationsObserved      *bool                         `db:"violations_observed"`
	Remarks                 *string                       `db:"remarks"`
	SubmittedAt             *time.Time                    `db:"submitted_at"`
	CreatedAt               time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time                     `json:"updated_at" db:"updated_at"`
}

type EvaluationReports []EvaluationReport

func (r *EvaluationReport) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return nil, nil
}

func (r *EvaluationReport) TableName() string {
	return "evaluation_reports"
}
