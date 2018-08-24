package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// ShipmentOffer maps a Transportation Service Provider to a shipment,
// indicating that the shipment has been offered to that TSP.
type ShipmentOffer struct {
	ID                              uuid.UUID `json:"id" db:"id"`
	CreatedAt                       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at" db:"updated_at"`
	ShipmentID                      uuid.UUID `json:"shipment_id" db:"shipment_id"`
	Shipment                        Shipment  `belongs_to:"shipments"`
	TransportationServiceProviderID uuid.UUID `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	AdministrativeShipment          bool      `json:"administrative_shipment" db:"administrative_shipment"`
	Accepted                        *bool     `json:"accepted" db:"accepted"`
	RejectionReason                 *string   `json:"rejection_reason" db:"rejection_reason"`
}

// String is not required by pop and may be deleted
func (a ShipmentOffer) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// ShipmentOffers is not required by pop and may be deleted
type ShipmentOffers []ShipmentOffer

// String is not required by pop and may be deleted
func (a ShipmentOffers) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

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
