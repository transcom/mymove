package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// PPMShipmentStatus represents the status of an order record's lifecycle
type PPMShipmentStatus string

const (
	// PPMShipmentStatusDRAFT captures enum value "DRAFT"
	PPMShipmentStatusDRAFT PPMShipmentStatus = "DRAFT"
	// PPMShipmentStatusSUBMITTED captures enum value "SUBMITTED"
	PPMShipmentStatusSUBMITTED PPMShipmentStatus = "SUBMITTED"
	// PPMShipmentStatusAPPROVED captures enum value "APPROVED"
	PPMShipmentStatusAPPROVED PPMShipmentStatus = "APPROVED"
	// PPMShipmentStatusPAYMENTREQUESTED captures enum value "PAYMENT_REQUESTED"
	PPMShipmentStatusPAYMENTREQUESTED PPMShipmentStatus = "PAYMENT_REQUESTED"
	// PPMShipmentStatusCOMPLETED captures enum value "COMPLETED"
	PPMShipmentStatusCOMPLETED PPMShipmentStatus = "COMPLETED"
	// PPMShipmentStatusCANCELED captures enum value "CANCELED"
	PPMShipmentStatusCANCELED PPMShipmentStatus = "CANCELED"
)

// PPMShipment is the portion of a move that a service member performs themselves
type PPMShipment struct {
	ID                             uuid.UUID         `json:"id" db:"id"`
	ShipmentID                     uuid.UUID         `json:"shipment_id" db:"shipment_id"` // Should this be MTOShipment?
	Shipment                       MTOShipment       `belongs_to:"shipment" fk_id:"shipment_id"`
	CreatedAt                      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time         `json:"updated_at" db:"updated_at"`
	Status                         PPMShipmentStatus `json:"status" db:"status"`
	ExpectedDepartureDate          *time.Time        `json:"expected_departure_date" db:"expected_departure_date"` // Originally OriginalMoveDate
	ActualMoveDate                 *time.Time        `json:"actual_move_date" db:"actual_move_date"`
	SubmitDate                     *time.Time        `json:"submit_date" db:"submit_date"`
	ReviewDate                     *time.Time        `json:"review_date" db:"review_date"`
	ApproveDate                    *time.Time        `json:"approve_date" db:"approve_date"`
	PickupPostalCode               *string           `json:"pickup_postal_code" db:"pickup_postal_code"`
	SecondaryPickupPostalCode      *string           `json:"secondary_pickup_postal_code" db:"secondary_pickup_postal_code"` // Originally AdditionalPickupPostalCode
	DestinationPostalCode          *string           `json:"destination_postal_code" db:"destination_postal_code"`
	SecondaryDestinationPostalCode *string           `json:"secondary_destination_postal_code" db:"secondary_destination_postal_code"`
	SitExpected                    *bool             `json:"sit_expected" db:"sit_expected"`         // Originally HasSit
	EstimatedWeight                *unit.Pound       `json:"estimated_weight" db:"estimated_weight"` // Originally WeightEstimate
	NetWeight                      *unit.Pound       `json:"net_weight" db:"net_weight"`
	HasProGear                     bool              `json:"has_pro_gear" db:"has_pro_gear"` // Can we get rid of this and just base it on if the pro gear weights are 0?
	ProGearWeight                  *int32            `json:"pro_gear_weight" db:"pro_gear_weight"`
	SpouseProGearWeight            *int32            `json:"spouse_pro_gear_weight" db:"spouse_pro_gear_weight"`
	EstimatedIncentive             *int32            `json:"estimated_incentive" db:"estimated_incentive"` // Originally IncentiveEstimate
	AdvanceRequested               bool              `json:"advance_requested" db:"advance_requested"`     // Originally HasRequestedAdvance
	AdvanceID                      *uuid.UUID        `json:"advance_id" db:"advance_id"`
	Advance                        *Reimbursement    `belongs_to:"reimbursements" fk_id:"advance_id"`
	AdvanceWorksheetID             *uuid.UUID        `json:"advance_worksheet_id" db:"advance_worksheet_id"`
	AdvanceWorksheet               Document          `belongs_to:"documents" fk_id:"advance_worksheet_id"`
}

// PPMShipments is a list of PPMs
type PPMShipments []PPMShipment
