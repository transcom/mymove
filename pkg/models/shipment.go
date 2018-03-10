package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
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
	CodeOfService             *string   `json:"code_of_service" db:"code_of_service"`
	Channel                   *string   `json:"channel" db:"channel"`
	GBLOC                     string    `json:"gbloc" db:"gbloc"`
	Market                    *string   `json:"market" db:"market"`
} // To do, Breanne: whatever's needed to integrate these three new fields throughout test data and tests and other functions. For instance, adding these new fields - possibly - to the PAS below. And adding it to the query that makes it.

// PossiblyAwardedShipment represents a single awarded shipment within a Service Member's move.
type PossiblyAwardedShipment struct {
	ID                              uuid.UUID  `db:"id"`
	CreatedAt                       time.Time  `db:"created_at"`
	UpdatedAt                       time.Time  `db:"updated_at"`
	TrafficDistributionListID       uuid.UUID  `db:"traffic_distribution_list_id"`
	PickupDate                      time.Time  `json:"pickup_date" db:"pickup_date"`
	TransportationServiceProviderID *uuid.UUID `db:"transportation_service_provider_id"`
	Accepted                        *bool      `json:"accepted" db:"accepted"`
	RejectionReason                 *string    `json:"rejection_reason" db:"rejection_reason"`
	AdministrativeShipment          *bool      `db:"administrative_shipment"`
}

// FetchPossiblyAwardedShipments runs the SQL query to fetch possibly awarded shipments from db
func FetchPossiblyAwardedShipments(dbConnection *pop.Connection) ([]PossiblyAwardedShipment, error) {
	shipments := []PossiblyAwardedShipment{}

	// TODO Can Q() be .All(&shipments)
	query := dbConnection.Q().LeftOuterJoin("shipment_awards", "shipment_awards.shipment_id=shipments.id")

	sql, args := query.ToSQL(&pop.Model{Value: Shipment{}},
		"shipments.id",
		"shipments.created_at",
		"shipments.updated_at",
		"shipments.traffic_distribution_list_id",
		"shipments.pickup_date",
		"shipment_awards.transportation_service_provider_id",
		"shipment_awards.administrative_shipment",
	)
	err := dbConnection.RawQuery(sql, args...).All(&shipments)
	return shipments, err
}

// FetchAwardedShipments looks up all unawarded shipments and returns them in the PossiblyAwardedShipment struct
// TODO: This is virtually identical to the function above, except it returns shipments that
//       are specifically awarded. Consolidate.
func FetchAwardedShipments(dbConnection *pop.Connection) ([]PossiblyAwardedShipment, error) {
	shipments := []PossiblyAwardedShipment{}

	sql := `SELECT
				shipments.id,
				shipments.traffic_distribution_list_id,
				shipment_awards.transportation_service_provider_id
			FROM shipments
			LEFT JOIN shipment_awards ON
				shipment_awards.shipment_id=shipments.id
			WHERE shipment_awards.id IS NULL`

	err := dbConnection.RawQuery(sql).All(&shipments)

	return shipments, err
}

// String is not required by pop and may be deleted
func (s Shipment) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Shipments is not required by pop and may be deleted
type Shipments []Shipment

// String is not required by pop and may be deleted
func (s Shipments) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s *Shipment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.TrafficDistributionListID, Name: "traffic_distribution_list_id"},
		&validators.StringIsPresent{Field: s.GBLOC, Name: "gbloc"},
	), nil
}
