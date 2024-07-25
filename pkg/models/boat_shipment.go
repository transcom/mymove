package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// BoatShipmentType represents the status of an order record's lifecycle
type BoatShipmentType string

const (
	// BoatShipmentTypeHaulAway captures enum value "HAUL_AWAY"
	BoatShipmentTypeHaulAway BoatShipmentType = "HAUL_AWAY"
	// BoatShipmentTypeTowAway captures enum value "TOW_AWAY"
	BoatShipmentTypeTowAway BoatShipmentType = "TOW_AWAY"
)

// AllowedBoatShipmentTypees is a list of all the allowed values for the Type of a BoatShipment as strings. Needed for
// validation.
var AllowedBoatShipmentTypes = []string{
	string(BoatShipmentTypeHaulAway),
	string(BoatShipmentTypeTowAway),
}

// BoatShipment is the portion of a move that a service member performs themselves
type BoatShipment struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	ShipmentID     uuid.UUID        `json:"shipment_id" db:"shipment_id"`
	Shipment       MTOShipment      `belongs_to:"mto_shipments" fk_id:"shipment_id"`
	Type           BoatShipmentType `json:"type" db:"type"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time       `json:"deleted_at" db:"deleted_at"`
	Year           int              `json:"year" db:"year"`
	Make           string           `json:"make" db:"make"`
	Model          string           `json:"model" db:"model"`
	LengthInInches int              `json:"length_in_inches" db:"length_in_inches"`
	WidthInInches  int              `json:"width_in_inches" db:"width_in_inches"`
	HeightInInches int              `json:"height_in_inches" db:"height_in_inches"`
	HasTrailer     bool             `json:"has_trailer" db:"has_trailer"`
	IsRoadworthy   bool             `json:"is_roadworthy" db:"is_roadworthy"`
}

// TableName overrides the table name used by Pop.
func (b BoatShipment) TableName() string {
	return "boat_shipments"
}

// BoatShipments is a list of Boats
type BoatShipments []BoatShipment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate,
// pop.ValidateAndUpdate) method. This should contain validation that is for data integrity. Business validation should
// occur in service objects.
func (b BoatShipment) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Name: "ShipmentID", Field: b.ShipmentID},
		&OptionalTimeIsPresent{Name: "DeletedAt", Field: b.DeletedAt},
		&validators.StringInclusion{Name: "Type", Field: string(b.Type), List: AllowedBoatShipmentTypes},
		&validators.IntIsGreaterThan{Name: "Year", Field: b.Year, Compared: 0},
		&validators.StringIsPresent{Name: "Make", Field: b.Make},
		&validators.StringIsPresent{Name: "Model", Field: b.Model},
		&validators.IntIsGreaterThan{Name: "LengthInInches", Field: b.LengthInInches, Compared: 0},
		&validators.IntIsGreaterThan{Name: "WidthInInches", Field: b.WidthInInches, Compared: 0},
		&validators.IntIsGreaterThan{Name: "HeightInInches", Field: b.HeightInInches, Compared: 0},
		&CannotBeTrueIfFalse{Field1: b.IsRoadworthy, Name1: "IsRoadworthy", Field2: b.HasTrailer, Name2: "HasTrailer"},
	), nil

}

// FetchBoatShipmentByBoatShipmentID returns a Boat Shipment for a given id
func FetchBoatShipmentByBoatShipmentID(db *pop.Connection, boatShipmentID uuid.UUID) (*BoatShipment, error) {
	var boatShipment BoatShipment
	err := db.Q().Find(&boatShipment, boatShipmentID)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}
	return &boatShipment, nil
}
