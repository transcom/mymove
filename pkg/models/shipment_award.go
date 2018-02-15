package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

// ShipmentAward maps a Transportation Service Provider to a shipment,
// indicating that the shipment has been awarded to that TSP.
type ShipmentAward struct {
	ID                              uuid.UUID `json:"id" db:"id"`
	CreatedAt                       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at" db:"updated_at"`
	ShipmentID                      uuid.UUID `json:"shipment_id" db:"shipment_id"`
	TransportationServiceProviderID uuid.UUID `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	AdministrativeShipment          bool      `json:"administrative_shipment" db:"administrative_shipment"`
}

// String is not required by pop and may be deleted
func (a ShipmentAward) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// ShipmentAwards is not required by pop and may be deleted
type ShipmentAwards []ShipmentAward

// String is not required by pop and may be deleted
func (a ShipmentAwards) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (a *ShipmentAward) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: a.ShipmentID, Name: "ShipmentID"},
		&validators.UUIDIsPresent{Field: a.TransportationServiceProviderID, Name: "TransportationServiceProviderID"},
	), nil
}

// CreateShipmentAward connects a shipment to a transportation service provider. This
// function assumes that the match has been validated by the caller.
func CreateShipmentAward(tx *pop.Connection,
	shipmentID uuid.UUID,
	tspID uuid.UUID,
	administrativeShipment bool) error {

	shipmentAward := ShipmentAward{
		ShipmentID:                      shipmentID,
		TransportationServiceProviderID: tspID,
		AdministrativeShipment:          administrativeShipment,
	}
	_, err := tx.ValidateAndSave(&shipmentAward)

	return err
}
