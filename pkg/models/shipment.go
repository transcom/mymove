package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
	"time"
)

// A Shipment represents a transfer of items form one location to another.
type Shipment struct {
	ID                              uuid.UUID     `json:"id" db:"id"`
	CreatedAt                       time.Time     `json:"created_at"`
	UpdatedAt                       time.Time     `json:"updated_at"`
	Name                            string        `json:"name"`
	PickupDate                      time.Time     `json:"pickup_date"`
	DeliveryDate                    time.Time     `json:"delivery_date"`
	TrafficDistributionListID       uuid.UUID     `json:"traffic_distribution_list_id"`
	TransportationServiceProviderID uuid.NullUUID `json:"transportation_service_provider_id"`
	AdministrativeShipment          bool
}

// Shipments is not required by pop and may be deleted
type Shipments []Shipment

// String is not required by pop and may be deleted
func (s Shipments) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Shipment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Shipment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Shipment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
