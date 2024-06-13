package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// GHCDieselFuelPrice represents the weekly national average diesel fuel price
type GHCDieselFuelPrice struct {
	ID                    uuid.UUID       `json:"id" db:"id"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
	FuelPriceInMillicents unit.Millicents `json:"fuel_price_in_millicents" db:"fuel_price_in_millicents"`
	PublicationDate       time.Time       `json:"publication_date" db:"publication_date"`
	EffectiveDate         time.Time       `json:"effective_date" db:"effective_date"`
	EndDate               time.Time       `json:"end_date" db:"end_date"`
}

// TableName overrides the table name used by Pop.
func (g GHCDieselFuelPrice) TableName() string {
	return "ghc_diesel_fuel_prices"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (g *GHCDieselFuelPrice) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: g.FuelPriceInMillicents.Int(), Name: "FuelPriceInMillicents"},
		&validators.TimeIsPresent{Field: g.PublicationDate, Name: "PublicationDate"},
	), nil
}
