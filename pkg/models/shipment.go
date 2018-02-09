package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

// Shipment represents a single shipment within a Service Member's move.
type Shipment struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
	PickupDate                time.Time `json:"pickup_date" db:"pickup_date"`
	DeliveryDate              time.Time `json:"delivery_date" db:"delivery_date"`
	TrafficDistributionListID uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
}

// PossiblyAwardedShipment represents a single awarded shipment within a Service Member's move.
type PossiblyAwardedShipment struct {
	ID                              uuid.UUID  `db:"id"`
	CreatedAt                       time.Time  `db:"created_at"`
	UpdatedAt                       time.Time  `db:"updated_at"`
	TrafficDistributionListID       uuid.UUID  `db:"traffic_distribution_list_id"`
	TransportationServiceProviderID *uuid.UUID `db:"transportation_service_provider_id"`
	AdministrativeShipment          *bool      `db:"administrative_shipment"`
}

// FetchPossiblyAwardedShipments runs the SQL query to fetch possibly awarded shipments from db
func FetchPossiblyAwardedShipments(dbConnection *pop.Connection) ([]PossiblyAwardedShipment, error) {
	shipments := []PossiblyAwardedShipment{}

	// TODO Can Q() be .All(&shipments)
	query := dbConnection.Q().LeftOuterJoin("awarded_shipments", "awarded_shipments.shipment_id=shipments.id")

	sql, args := query.ToSQL(&pop.Model{Value: Shipment{}},
		"shipments.id",
		"shipments.created_at",
		"shipments.updated_at",
		"shipments.traffic_distribution_list_id",
		"awarded_shipments.transportation_service_provider_id",
		"awarded_shipments.administrative_shipment",
	)
	err := dbConnection.RawQuery(sql, args...).All(&shipments)
	return shipments, err
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
