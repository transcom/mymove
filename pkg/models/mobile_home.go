package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type MobileHome struct {
	ID                             uuid.UUID   `json:"id" db:"id"`
	ShipmentID                     uuid.UUID   `json:"shipment_id" db:"shipment_id"`
	Shipment                       MTOShipment `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	Make                           string
	Model                          string    `db:"model"`
	Year                           int       `db:"mh_year"`
	Length                         int       `db:"mh_length"`
	Height                         int       `db:"height"`
	Width                          int       `db:"width"`
	RequestedPickupDate            time.Time `db:"requested_pickup_date"`
	RequestedDeliveryDate          time.Time `db:"requested_delivery_date"`
	PickupAddress                  string    `db:"pickup_address"`
	DestinationAddress             string    `db:"destination_address"`
	OriginAddress                  string    `db:"origin_address"`
	UpdatedAt                      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                      time.Time `json:"deleted_at" db:"deleted_at"`
	HasSecondaryPickupAddress      bool      `db:"has_secondary_pickup_address"`
	HasSecondaryDestinationAddress bool      `db:"has_secondary_destination_address"`
	SecondaryPickupAddress         string    `db:"secondary_pickup_address"`
	SecondaryDestinationAddress    string    `db:"secondary_destination_address"`
	ReceivingAgent                 string    `db:"receiving_agent"`
	CounselorRemarks               string    `db:"counselor_remarks"`
	CustomerRemarks                string    `db:"customer_remarks"`
}

// TableName overrides the table name used by Pop.
func (mh MobileHome) TableName() string {
	return "mobile_home"
}

// A list of Mobile homes
type MobileHomes []MobileHome

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
//TODO: KOSEY only add *fields to be validated
func (mh MobileHome) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "ShipmentID", Field: mh.ShipmentID},
	), nil

}
