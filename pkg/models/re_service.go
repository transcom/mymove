package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type ReServiceName string

const (
	DomesticLinehaul           ReServiceName = "Dom. Linehaul"
	FuelSurcharge              ReServiceName = "Fuel Surcharge"
	DomesticOriginPrice        ReServiceName = "Dom. Origin Price"
	DomesticDestinationPrice   ReServiceName = "Dom. Destination Price"
	DomesticPacking            ReServiceName = "Dom. Packing"
	DomesticUnpacking          ReServiceName = "Dom. Unpacking"
	DomesticShorthaul          ReServiceName = "Dom. Shorthaul"
	DomesticNTSPackingFactor   ReServiceName = "Dom. NTS Packing Factor"
	DomesticMobileHomeFactor   ReServiceName = "Dom. Mobile Home Factor"
	DomesticHaulAwayBoatFactor ReServiceName = "Dom. Haul Away Boat Factor"
	DomesticTowAwayBoatFactor  ReServiceName = "Dom. Tow Away Boat Factor"
)

// ReService model struct
type ReService struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReServices is not required by pop and may be deleted
type ReServices []ReService

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReService) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Code, Name: "Code"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}
