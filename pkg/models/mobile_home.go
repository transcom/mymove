package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type MobileHome struct {
	ID                        uuid.UUID   `json:"id" db:"id"`
	ShipmentID                uuid.UUID   `json:"shipment_id" db:"shipment_id"`
	Shipment                  MTOShipment `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	Make                      *string
	Model                     *string    `db:"model"`
	Year                      *int       `db:"year"`
	Length                    *int       `db:"length"`
	Width                     *int       `db:"width"`
	Height                    *int       `db:"height"`
	UpdatedAt                 time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at" db:"deleted_at"`
	RequestedPickupDate       *time.Time `db:"requested_pickup_date"`
	PickupLocation            *string    `db:"pickup_location"`
	RequestedDeliveryDate     *time.Time `db:"requested_delivery_date"`
	Dimensions                *string    `db:"dimensions"`
	OrginAddress              *string    `db:"origin_address"`
	HasSecondaryPickupAddress *bool      `db:"has_secondary_pickup_address"`
	SecondaryPickupAddress    *string    `db:"secondary_pickup_address"`
	ReceivingAgent            *string    `db:"receiving_agent"`
	CounselorRemarks          *string    `db:"counselor_remarks"`
	CustomerRemarks           *string    `db:"customer_remarks"`
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
func (mh MobileHome) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "ShipmentID", Field: mh.ShipmentID},
		&OptionalTimeIsPresent{Name: "DeletedAt", Field: mh.DeletedAt},
		&OptionalIntIsPositive{Name: "Height", Field: mh.Height},

			), nil

}