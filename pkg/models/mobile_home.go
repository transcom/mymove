package models

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type MobileHome struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	ShipmentID     uuid.UUID   `json:"shipment_id" db:"shipment_id"`
	Shipment       MTOShipment `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	Make           string      `json:"make" db:"make"`
	Model          string      `json:"model" db:"model"`
	Year           int         `json:"year" db:"year"`
	LengthInInches *int        `json:"length_in_inches" db:"length_in_inches"`
	HeightInInches *int        `json:"height_in_inches" db:"height_in_inches"`
	WidthInInches  *int        `json:"width_in_inches" db:"width_in_inches"`
}

// TableName overrides the table name used by Pop.
func (mh MobileHome) TableName() string {
	return "mobile_homes"
}

// A list of Mobile homes
type MobileHomes []MobileHome

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (mh MobileHome) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "ShipmentID", Field: mh.ShipmentID},
		&validators.StringIsPresent{Name: "Make", Field: mh.Make},
		&validators.StringIsPresent{Name: "Model", Field: mh.Model},
		&validators.IntIsGreaterThan{Name: "Year", Field: mh.Year, Compared: 0},
		&validators.IntIsGreaterThan{Name: "Length", Field: *mh.LengthInInches, Compared: 0},
		&validators.IntIsGreaterThan{Name: "Height", Field: *mh.HeightInInches, Compared: 0},
		&validators.IntIsGreaterThan{Name: "Width", Field: *mh.WidthInInches, Compared: 0},
	), nil
}

// Returns a Mobile Home Shipment for a given id
func FetchMobileHomeShipmentByMobileHomeShipmentID(db *pop.Connection, mobileHomeShipmentID uuid.UUID) (*MobileHome, error) {
	var mobileHome MobileHome
	err := db.Q().Find(&mobileHome, mobileHomeShipmentID)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}
	return &mobileHome, nil
}
