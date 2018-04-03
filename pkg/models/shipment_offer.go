package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// ShipmentOffer maps a Transportation Service Provider to a shipment,
// indicating that the shipment has been offered to that TSP.
type ShipmentOffer struct {
	ID                              uuid.UUID `db:"id"`
	CreatedAt                       time.Time `db:"created_at"`
	UpdatedAt                       time.Time `db:"updated_at"`
	ShipmentID                      uuid.UUID `db:"shipment_id"`
	TransportationServiceProviderID uuid.UUID `db:"transportation_service_provider_id"`
	AdministrativeShipment          bool      `db:"administrative_shipment"`
	Accepted                        *bool     `db:"accepted"`
	RejectionReason                 *string   `db:"rejection_reason"`
}

// ShipmentOffers is not required by pop and may be deleted
type ShipmentOffers []ShipmentOffer

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (a *ShipmentOffer) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: a.ShipmentID, Name: "ShipmentID"},
		&validators.UUIDIsPresent{Field: a.TransportationServiceProviderID, Name: "TransportationServiceProviderID"},
	), nil
}

// CreateShipmentOffer connects a shipment to a transportation service provider. This
// function assumes that the match has been validated by the caller.
func CreateShipmentOffer(tx *pop.Connection,
	shipmentID uuid.UUID,
	tspID uuid.UUID,
	administrativeShipment bool) (*ShipmentOffer, error) {

	shipmentOffer := ShipmentOffer{
		ShipmentID:                      shipmentID,
		TransportationServiceProviderID: tspID,
		AdministrativeShipment:          administrativeShipment,
	}
	_, err := tx.ValidateAndSave(&shipmentOffer)

	return &shipmentOffer, err
}
