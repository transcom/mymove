package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// MTOShipmentType represents the type of shipments the mto shipment is
type MTOShipmentType string

// using these also in move.go selected move type
const (
	// NTSRaw is the raw string value of the NTS Shipment Type
	NTSRaw = "HHG_INTO_NTS_DOMESTIC"
	// NTSrRaw is the raw string value of the NTSr Shipment Type
	NTSrRaw = "HHG_OUTOF_NTS_DOMESTIC"
)

const (
	// MTOShipmentTypeHHG is an HHG Shipment Type default
	MTOShipmentTypeHHG MTOShipmentType = "HHG"
	// MTOShipmentTypeInternationalHHG is a Shipment Type for International HHG
	MTOShipmentTypeInternationalHHG MTOShipmentType = "INTERNATIONAL_HHG"
	// MTOShipmentTypeInternationalUB is a Shipment Type for International UB
	MTOShipmentTypeInternationalUB MTOShipmentType = "INTERNATIONAL_UB"
	// MTOShipmentTypeHHGLongHaulDom is an HHG Shipment Type for Longhaul Domestic
	MTOShipmentTypeHHGLongHaulDom MTOShipmentType = "HHG_LONGHAUL_DOMESTIC"
	// MTOShipmentTypeHHGShortHaulDom is an HHG Shipment Type for Shothaul Domestic
	MTOShipmentTypeHHGShortHaulDom MTOShipmentType = "HHG_SHORTHAUL_DOMESTIC"
	// MTOShipmentTypeHHGIntoNTSDom is an HHG Shipment Type for going into NTS Domestic
	MTOShipmentTypeHHGIntoNTSDom MTOShipmentType = NTSRaw
	// MTOShipmentTypeHHGOutOfNTSDom is an HHG Shipment Type for going out of NTS Domestic
	MTOShipmentTypeHHGOutOfNTSDom MTOShipmentType = NTSrRaw
	// MTOShipmentTypeMotorhome is a Shipment Type for Motorhome
	MTOShipmentTypeMotorhome MTOShipmentType = "MOTORHOME"
	// MTOShipmentTypeBoatHaulAway is a Shipment Type for Boat Haul Away
	MTOShipmentTypeBoatHaulAway MTOShipmentType = "BOAT_HAUL_AWAY"
	// MTOShipmentTypeBoatTowAway is a Shipment Type for Boat Tow Away
	MTOShipmentTypeBoatTowAway MTOShipmentType = "BOAT_TOW_AWAY"
)

// MTOShipmentStatus represents the possible statuses for a mto shipment
type MTOShipmentStatus string

const (
	// MTOShipmentStatusDraft is the draft status type for MTO Shipments
	MTOShipmentStatusDraft MTOShipmentStatus = "DRAFT"
	// MTOShipmentStatusSubmitted is the submitted status type for MTO Shipments
	MTOShipmentStatusSubmitted MTOShipmentStatus = "SUBMITTED"
	// MTOShipmentStatusApproved is the approved status type for MTO Shipments
	MTOShipmentStatusApproved MTOShipmentStatus = "APPROVED"
	// MTOShipmentStatusRejected is the rejected status type for MTO Shipments
	MTOShipmentStatusRejected MTOShipmentStatus = "REJECTED"
)

// MTOShipment is an object representing data for a move task order shipment
type MTOShipment struct {
	ID                               uuid.UUID         `db:"id"`
	MoveTaskOrder                    Move              `belongs_to:"moves"`
	MoveTaskOrderID                  uuid.UUID         `db:"move_id"`
	ScheduledPickupDate              *time.Time        `db:"scheduled_pickup_date"`
	RequestedPickupDate              *time.Time        `db:"requested_pickup_date"`
	RequestedDeliveryDate            *time.Time        `db:"requested_delivery_date"`
	ApprovedDate                     *time.Time        `db:"approved_date"`
	FirstAvailableDeliveryDate       *time.Time        `db:"first_available_delivery_date"`
	ActualPickupDate                 *time.Time        `db:"actual_pickup_date"`
	RequiredDeliveryDate             *time.Time        `db:"required_delivery_date"`
	CustomerRemarks                  *string           `db:"customer_remarks"`
	PickupAddress                    *Address          `belongs_to:"addresses"`
	PickupAddressID                  *uuid.UUID        `db:"pickup_address_id"`
	DestinationAddress               *Address          `belongs_to:"addresses"`
	DestinationAddressID             *uuid.UUID        `db:"destination_address_id"`
	MTOAgents                        MTOAgents         `has_many:"mto_agents" fk_id:"mto_shipment_id"`
	MTOServiceItems                  MTOServiceItems   `has_many:"mto_service_items" fk_id:"mto_shipment_id"`
	SecondaryPickupAddress           *Address          `belongs_to:"addresses"`
	SecondaryPickupAddressID         *uuid.UUID        `db:"secondary_pickup_address_id"`
	SecondaryDeliveryAddress         *Address          `belongs_to:"addresses"`
	SecondaryDeliveryAddressID       *uuid.UUID        `db:"secondary_delivery_address_id"`
	PrimeEstimatedWeight             *unit.Pound       `db:"prime_estimated_weight"`
	PrimeEstimatedWeightRecordedDate *time.Time        `db:"prime_estimated_weight_recorded_date"`
	PrimeActualWeight                *unit.Pound       `db:"prime_actual_weight"`
	ShipmentType                     MTOShipmentType   `db:"shipment_type"`
	Status                           MTOShipmentStatus `db:"status"`
	RejectionReason                  *string           `db:"rejection_reason"`
	Distance                         *unit.Miles       `db:"distance"`
	CreatedAt                        time.Time         `db:"created_at"`
	UpdatedAt                        time.Time         `db:"updated_at"`
}

// MTOShipments is a list of mto shipments
type MTOShipments []MTOShipment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *MTOShipment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringInclusion{Field: string(m.Status), Name: "Status", List: []string{
		string(MTOShipmentStatusApproved),
		string(MTOShipmentStatusRejected),
		string(MTOShipmentStatusSubmitted),
		string(MTOShipmentStatusDraft),
	}})
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MoveTaskOrderID, Name: "MoveTaskOrderID"})
	if m.PrimeEstimatedWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.PrimeEstimatedWeight.Int(), Compared: -1, Name: "PrimeEstimatedWeight"})
	}
	if m.PrimeActualWeight != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: m.PrimeActualWeight.Int(), Compared: -1, Name: "PrimeActualWeight"})
	}
	if m.Status == MTOShipmentStatusRejected {
		var rejectionReason string
		if m.RejectionReason != nil {
			rejectionReason = *m.RejectionReason
		}
		vs = append(vs, &validators.StringIsPresent{Field: rejectionReason, Name: "RejectionReason"})
	}
	return validate.Validate(vs...), nil
}

// TableName overrides the table name used by Pop.
func (m MTOShipment) TableName() string {
	return "mto_shipments"
}
