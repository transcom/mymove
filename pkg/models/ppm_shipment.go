package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// PPMShipmentStatus represents the status of an order record's lifecycle
type PPMShipmentStatus string

const (
	// PPMShipmentStatusDraft captures enum value "DRAFT"
	PPMShipmentStatusDraft PPMShipmentStatus = "DRAFT"
	// PPMShipmentStatusSubmitted captures enum value "SUBMITTED"
	PPMShipmentStatusSubmitted PPMShipmentStatus = "SUBMITTED"
	// PPMShipmentStatusWaitingOnCustomer captures enum value "WAITING_ON_CUSTOMER"
	PPMShipmentStatusWaitingOnCustomer PPMShipmentStatus = "WAITING_ON_CUSTOMER"
	// PPMShipmentStatusNeedsAdvanceApproval captures enum value "NEEDS_ADVANCE_APPROVAL"
	PPMShipmentStatusNeedsAdvanceApproval PPMShipmentStatus = "NEEDS_ADVANCE_APPROVAL"
	// PPMShipmentStatusNeedsPaymentApproval captures enum value "NEEDS_PAYMENT_APPROVAL"
	PPMShipmentStatusNeedsPaymentApproval PPMShipmentStatus = "NEEDS_PAYMENT_APPROVAL"
	// PPMShipmentStatusPaymentApproved captures enum value "PAYMENT_APPROVED"
	PPMShipmentStatusPaymentApproved PPMShipmentStatus = "PAYMENT_APPROVED"
)

// PPMShipment is the portion of a move that a service member performs themselves
type PPMShipment struct {
	ID                             uuid.UUID         `json:"id" db:"id"`
	ShipmentID                     uuid.UUID         `json:"shipment_id" db:"shipment_id"`
	Shipment                       MTOShipment       `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	CreatedAt                      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time         `json:"updated_at" db:"updated_at"`
	Status                         PPMShipmentStatus `json:"status" db:"status"`
	ExpectedDepartureDate          time.Time         `json:"expected_departure_date" db:"expected_departure_date"`
	ActualMoveDate                 *time.Time        `json:"actual_move_date" db:"actual_move_date"`
	SubmittedAt                    *time.Time        `json:"submitted_at" db:"submitted_at"`
	ReviewedAt                     *time.Time        `json:"reviewed_at" db:"reviewed_at"`
	ApprovedAt                     *time.Time        `json:"approved_at" db:"approved_at"`
	PickupPostalCode               string            `json:"pickup_postal_code" db:"pickup_postal_code"`
	SecondaryPickupPostalCode      *string           `json:"secondary_pickup_postal_code" db:"secondary_pickup_postal_code"`
	DestinationPostalCode          string            `json:"destination_postal_code" db:"destination_postal_code"`
	SecondaryDestinationPostalCode *string           `json:"secondary_destination_postal_code" db:"secondary_destination_postal_code"`
	SitExpected                    bool              `json:"sit_expected" db:"sit_expected"`
	EstimatedWeight                *unit.Pound       `json:"estimated_weight" db:"estimated_weight"`
	NetWeight                      *unit.Pound       `json:"net_weight" db:"net_weight"`
	HasProGear                     *bool             `json:"has_pro_gear" db:"has_pro_gear"`
	ProGearWeight                  *unit.Pound       `json:"pro_gear_weight" db:"pro_gear_weight"`
	SpouseProGearWeight            *unit.Pound       `json:"spouse_pro_gear_weight" db:"spouse_pro_gear_weight"`
	EstimatedIncentive             *int32            `json:"estimated_incentive" db:"estimated_incentive"`
	Advance                        *unit.Cents       `json:"advance" db:"advance"`
}

// PPMShipments is a list of PPMs
type PPMShipments []PPMShipment

// TableName overrides the table name used by Pop. By default it tries using the name `ppmshipments`.
func (p PPMShipment) TableName() string {
	return "ppm_shipments"
}
