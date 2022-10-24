package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type EvaluationReportType string

const (
	EvaluationReportTypeShipment   EvaluationReportType = "SHIPMENT"
	EvaluationReportTypeCounseling EvaluationReportType = "COUNSELING"
)

type EvaluationReportInspectionType string

const (
	EvaluationReportInspectionTypeDataReview EvaluationReportInspectionType = "DATA_REVIEW"
	EvaluationReportInspectionTypePhysical   EvaluationReportInspectionType = "PHYSICAL"
	EvaluationReportInspectionTypeVirtual    EvaluationReportInspectionType = "VIRTUAL"
)

type EvaluationReportLocationType string

const (
	EvaluationReportLocationTypeOrigin      EvaluationReportLocationType = "ORIGIN"
	EvaluationReportLocationTypeDestination EvaluationReportLocationType = "DESTINATION"
	EvaluationReportLocationTypeOther       EvaluationReportLocationType = "OTHER"
)

type EvaluationReport struct {
	ID                                 uuid.UUID                       `json:"id" db:"id"`
	OfficeUser                         OfficeUser                      `belongs_to:"office_users" fk_id:"office_user_id"`
	OfficeUserID                       uuid.UUID                       `db:"office_user_id"`
	Move                               Move                            `belongs_to:"moves" fk_id:"move_id"`
	MoveID                             uuid.UUID                       `db:"move_id"`
	Shipment                           *MTOShipment                    `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	ShipmentID                         *uuid.UUID                      `json:"shipment_id" db:"shipment_id"`
	Type                               EvaluationReportType            `json:"type" db:"type"`
	InspectionDate                     *time.Time                      `json:"inspection_date" db:"inspection_date"`
	InspectionType                     *EvaluationReportInspectionType `json:"inspection_type" db:"inspection_type"`
	TravelTimeMinutes                  *int                            `json:"travel_time_minutes" db:"travel_time_minutes"`
	Location                           *EvaluationReportLocationType   `json:"location" db:"location"`
	LocationDescription                *string                         `json:"location_description" db:"location_description"`
	ObservedShipmentDeliveryDate       *time.Time                      `json:"observed_shipment_delivery_date" db:"observed_shipment_delivery_date"`
	ObservedShipmentPhysicalPickupDate *time.Time                      `json:"observed_shipment_physical_pickup_date" db:"observed_shipment_physical_pickup_date"`
	EvaluationLengthMinutes            *int                            `json:"evaluation_length_minutes" db:"evaluation_length_minutes"`
	ViolationsObserved                 *bool                           `json:"violations_observed" db:"violations_observed"`
	Remarks                            *string                         `json:"remarks" db:"remarks"`
	SeriousIncident                    *bool                           `json:"serious_incident" db:"serious_incident"`
	SeriousIncidentDesc                *string                         `json:"serious_incident_desc" db:"serious_incident_desc"`
	ObservedClaimsResponseDate         *time.Time                      `json:"observed_claims_response_date" db:"observed_claims_response_date"`
	ObservedPickupDate                 *time.Time                      `json:"observed_pickup_date" db:"observed_pickup_date"`
	ObservedPickupSpreadStartDate      *time.Time                      `json:"observed_pickup_spread_start_date" db:"observed_pickup_spread_start_date"`
	ObservedPickupSpreadEndDate        *time.Time                      `json:"observed_pickup_spread_end_date" db:"observed_pickup_spread_end_date"`
	ObservedDeliveryDate               *time.Time                      `json:"observed_delivery_date" db:"observed_delivery_date"`
	SubmittedAt                        *time.Time                      `json:"submitted_at" db:"submitted_at"`
	DeletedAt                          *time.Time                      `db:"deleted_at"`
	CreatedAt                          time.Time                       `json:"created_at" db:"created_at"`
	UpdatedAt                          time.Time                       `json:"updated_at" db:"updated_at"`
	ReportViolations                   ReportViolations                `json:"report_violation,omitempty" fk_id:"report_id" has_many:"report_violation"`
}

// EvaluationReports is not required by pop and may be deleted
type EvaluationReports []EvaluationReport

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *EvaluationReport) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	if r.ShipmentID != nil {
		vs = append(vs, &validators.StringsMatch{Name: "Type", Field: string(r.Type), Field2: string(EvaluationReportTypeShipment)})
	}
	if r.TravelTimeMinutes != nil {
		vs = append(vs, &validators.StringsMatch{Field: string(*r.InspectionType), Name: "InspectionType", Field2: string(EvaluationReportInspectionTypePhysical)})
	}

	if r.ObservedShipmentDeliveryDate != nil {
		vs = append(vs, &validators.StringsMatch{Field: string(*r.InspectionType), Name: "InspectionType", Field2: string(EvaluationReportInspectionTypePhysical)})
	}

	if r.ObservedShipmentPhysicalPickupDate != nil {
		vs = append(vs, &validators.StringsMatch{Field: string(*r.InspectionType), Name: "InspectionType", Field2: string(EvaluationReportInspectionTypePhysical)})
	}

	if r.LocationDescription != nil {
		vs = append(vs, &validators.StringsMatch{Field: string(*r.Location), Name: "Location", Field2: string(EvaluationReportLocationTypeOther)})
	}
	verrs := validate.Validate(vs...)
	if r.Type == EvaluationReportTypeShipment && r.ShipmentID == nil {
		verrs.Add(validators.GenerateKey("ShipmentID"), "If report type is SHIPMENT, ShipmentID must not be null")
	}
	return verrs, nil
}

func (r *EvaluationReport) TableName() string {
	return "evaluation_reports"
}
