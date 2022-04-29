package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// ReweighRequester is the actor who initiated the reweigh request for the shipment
type ReweighRequester string

const (
	// ReweighRequesterCustomer represents the customer requesting a reweigh
	ReweighRequesterCustomer ReweighRequester = "CUSTOMER"
	// ReweighRequesterPrime represents the prime mover requesting a reweigh
	ReweighRequesterPrime ReweighRequester = "PRIME"
	// ReweighRequesterSystem represents the milmove system triggering a reweigh
	ReweighRequesterSystem ReweighRequester = "SYSTEM"
	// ReweighRequesterTOO represents the TOO office user requesting a reweigh
	ReweighRequesterTOO ReweighRequester = "TOO"
)

var requestedByValues = []string{
	string(ReweighRequesterCustomer),
	string(ReweighRequesterPrime),
	string(ReweighRequesterSystem),
	string(ReweighRequesterTOO),
}

// Reweigh represents a request for the prime mover to reweigh a shipment or provide verification why they could not
type Reweigh struct {
	ID                     uuid.UUID        `db:"id"`
	RequestedAt            time.Time        `db:"requested_at"`
	RequestedBy            ReweighRequester `db:"requested_by"`
	Shipment               MTOShipment      `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	ShipmentID             uuid.UUID        `db:"shipment_id"`
	VerificationProvidedAt *time.Time       `db:"verification_provided_at"`
	VerificationReason     *string          `db:"verification_reason"`
	Weight                 *unit.Pound      `db:"weight"`
	CreatedAt              time.Time        `db:"created_at"`
	UpdatedAt              time.Time        `db:"updated_at"`
}

// Validate ensures the reweigh fields have the required and optional valid values prior to saving
func (r *Reweigh) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: r.RequestedAt, Name: "RequestedAt"},
		&validators.StringInclusion{Field: string(r.RequestedBy), Name: "RequestedBy", List: requestedByValues},
		&validators.UUIDIsPresent{Field: r.ShipmentID, Name: "ShipmentID"},
		&OptionalTimeIsPresent{Field: r.VerificationProvidedAt, Name: "VerificationProvidedAt"},
		&StringIsNilOrNotBlank{Field: r.VerificationReason, Name: "VerificationReason"},
		&OptionalPoundIsPositive{Field: r.Weight, Name: "Weight"},
	), nil
}
