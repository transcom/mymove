package models

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/db/utilities"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
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

// PPMAdvanceStatus represents the status of an advance that can be approved, edited or rejected by a SC
type PPMAdvanceStatus string

const (
	// PPMAdvanceStatusApproved captures enum value "APPROVED"
	PPMAdvanceStatusApproved PPMAdvanceStatus = "APPROVED"
	// PPMAdvanceStatusEdited captures enum value "EDITED"
	PPMAdvanceStatusEdited PPMAdvanceStatus = "EDITED"
	// PPMAdvanceStatusRejected captures enum value "REJECTED"
	PPMAdvanceStatusRejected PPMAdvanceStatus = "REJECTED"
)

// SITLocationType represents whether the SIT at the origin or destination
type SITLocationType string

const (
	// SITLocationTypeOrigin captures enum value "ORIGIN"
	SITLocationTypeOrigin SITLocationType = "ORIGIN"
	// SITLocationTypeDestination captures enum value "DESTINATION"
	SITLocationTypeDestination SITLocationType = "DESTINATION"
)

// PPMDocumentStatus represents the status of a PPMShipment's documents. Lives here since we have multiple PPM document
// models.
type PPMDocumentStatus string

const (
	// PPMDocumentStatusApproved captures enum value "APPROVED"
	PPMDocumentStatusApproved PPMDocumentStatus = "APPROVED"
	// PPMDocumentStatusExcluded captures enum value "EXCLUDED"
	PPMDocumentStatusExcluded PPMDocumentStatus = "EXCLUDED"
	// PPMDocumentStatusRejected captures enum value "REJECTED"
	PPMDocumentStatusRejected PPMDocumentStatus = "REJECTED"
)

var AllowedPPMDocumentStatuses = []string{
	string(PPMDocumentStatusApproved),
	string(PPMDocumentStatusExcluded),
	string(PPMDocumentStatusRejected),
}

// PPMShipment is the portion of a move that a service member performs themselves
type PPMShipment struct {
	ID                             uuid.UUID         `json:"id" db:"id"`
	ShipmentID                     uuid.UUID         `json:"shipment_id" db:"shipment_id"`
	Shipment                       MTOShipment       `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	CreatedAt                      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time         `json:"updated_at" db:"updated_at"`
	DeletedAt                      *time.Time        `json:"deleted_at" db:"deleted_at"`
	Status                         PPMShipmentStatus `json:"status" db:"status"`
	ExpectedDepartureDate          time.Time         `json:"expected_departure_date" db:"expected_departure_date"`
	ActualMoveDate                 *time.Time        `json:"actual_move_date" db:"actual_move_date"`
	SubmittedAt                    *time.Time        `json:"submitted_at" db:"submitted_at"`
	ReviewedAt                     *time.Time        `json:"reviewed_at" db:"reviewed_at"`
	ApprovedAt                     *time.Time        `json:"approved_at" db:"approved_at"`
	PickupPostalCode               string            `json:"pickup_postal_code" db:"pickup_postal_code"`
	SecondaryPickupPostalCode      *string           `json:"secondary_pickup_postal_code" db:"secondary_pickup_postal_code"`
	ActualPickupPostalCode         *string           `json:"actual_pickup_postal_code" db:"actual_pickup_postal_code"`
	DestinationPostalCode          string            `json:"destination_postal_code" db:"destination_postal_code"`
	SecondaryDestinationPostalCode *string           `json:"secondary_destination_postal_code" db:"secondary_destination_postal_code"`
	ActualDestinationPostalCode    *string           `json:"actual_destination_postal_code" db:"actual_destination_postal_code"`
	EstimatedWeight                *unit.Pound       `json:"estimated_weight" db:"estimated_weight"`
	NetWeight                      *unit.Pound       `json:"net_weight" db:"net_weight"`
	HasProGear                     *bool             `json:"has_pro_gear" db:"has_pro_gear"`
	ProGearWeight                  *unit.Pound       `json:"pro_gear_weight" db:"pro_gear_weight"`
	SpouseProGearWeight            *unit.Pound       `json:"spouse_pro_gear_weight" db:"spouse_pro_gear_weight"`
	EstimatedIncentive             *unit.Cents       `json:"estimated_incentive" db:"estimated_incentive"`
	HasRequestedAdvance            *bool             `json:"has_requested_advance" db:"has_requested_advance"`
	AdvanceAmountRequested         *unit.Cents       `json:"advance_amount_requested" db:"advance_amount_requested"`
	HasReceivedAdvance             *bool             `json:"has_received_advance" db:"has_received_advance"`
	AdvanceStatus                  *PPMAdvanceStatus `json:"advance_status" db:"advance_status"`
	AdvanceAmountReceived          *unit.Cents       `json:"advance_amount_received" db:"advance_amount_received"`
	SITExpected                    *bool             `json:"sit_expected" db:"sit_expected"`
	SITLocation                    *SITLocationType  `json:"sit_location" db:"sit_location"`
	SITEstimatedWeight             *unit.Pound       `json:"sit_estimated_weight" db:"sit_estimated_weight"`
	SITEstimatedEntryDate          *time.Time        `json:"sit_estimated_entry_date" db:"sit_estimated_entry_date"`
	SITEstimatedDepartureDate      *time.Time        `json:"sit_estimated_departure_date" db:"sit_estimated_departure_date"`
	SITEstimatedCost               *unit.Cents       `json:"sit_estimated_cost" db:"sit_estimated_cost"`
	WeightTickets                  WeightTickets     `has_many:"weight_tickets" fk_id:"ppm_shipment_id" order_by:"created_at asc"`
	MovingExpenses                 MovingExpenses    `has_many:"moving_expenses" fk_id:"ppm_shipment_id" order_by:"created_at asc"`
}

// PPMShipments is a list of PPMs
type PPMShipments []PPMShipment

// TableName overrides the table name used by Pop. By default it tries using the name `ppmshipments`.
func (p PPMShipment) TableName() string {
	return "ppm_shipments"
}

func FetchPPMShipmentFromMTOShipmentID(db *pop.Connection, mtoShipmentID uuid.UUID) (*PPMShipment, error) {
	var ppmShipment PPMShipment

	err := db.Scope(utilities.ExcludeDeletedScope()).EagerPreload("Shipment").
		Where("ppm_shipments.shipment_id = ?", mtoShipmentID).
		First(&ppmShipment)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "")
		}
	}
	return &ppmShipment, nil
}
