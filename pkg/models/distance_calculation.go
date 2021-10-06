package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// See: pkg/route/planner.go for more info on this interface
type distanceCalculator interface {
	Zip5TransitDistanceLineHaul(string, string) (int, error)
	TransitDistance(*Address, *Address) (int, error)
}

// DistanceCalculation represents a distance calculation in miles between an origin and destination address
type DistanceCalculation struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	OriginAddressID      uuid.UUID `json:"origin_address_id" db:"origin_address_id"`
	OriginAddress        Address   `belongs_to:"address" fk_id:"origin_address_id"`
	DestinationAddressID uuid.UUID `json:"destination_address_id" db:"destination_address_id"`
	DestinationAddress   Address   `belongs_to:"address" fk_id:"destination_address_id"`
	DistanceMiles        int       `json:"distance_miles" db:"distance_miles"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// NewDistanceCalculation performs a distance calculation and returns the resulting DistanceCalculation model
func NewDistanceCalculation(planner distanceCalculator, origin Address, destination Address, useZipOnly bool) (DistanceCalculation, error) {

	var distanceMiles int
	var err error

	if useZipOnly {
		distanceMiles, err = planner.Zip5TransitDistanceLineHaul(origin.PostalCode, destination.PostalCode)
	} else {
		distanceMiles, err = planner.TransitDistance(&origin, &destination)
	}
	if err != nil {
		return DistanceCalculation{}, err
	}

	distModel := DistanceCalculation{
		OriginAddress:        origin,
		OriginAddressID:      origin.ID,
		DestinationAddress:   destination,
		DestinationAddressID: destination.ID,
		DistanceMiles:        distanceMiles,
	}

	return distModel, nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *DistanceCalculation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: d.OriginAddressID, Name: "OriginAddressID"},
		&validators.UUIDIsPresent{Field: d.DestinationAddressID, Name: "DestinationAddressID"},
		&validators.IntIsPresent{Field: d.DistanceMiles, Name: "DistanceMiles"},
	), nil
}
