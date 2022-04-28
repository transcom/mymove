package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// ReDomesticLinehaulPrice represents a domestic linehaul price based on weight, mileage, etc.
type ReDomesticLinehaulPrice struct {
	ID                    uuid.UUID       `json:"id" db:"id"`
	ContractID            uuid.UUID       `json:"contract_id" db:"contract_id"`
	WeightLower           unit.Pound      `json:"weight_lower" db:"weight_lower"`
	WeightUpper           unit.Pound      `json:"weight_upper" db:"weight_upper"`
	MilesLower            int             `json:"miles_lower" db:"miles_lower"`
	MilesUpper            int             `json:"miles_upper" db:"miles_upper"`
	IsPeakPeriod          bool            `json:"is_peak_period" db:"is_peak_period"`
	DomesticServiceAreaID uuid.UUID       `json:"domestic_service_area_id" db:"domestic_service_area_id"`
	PriceMillicents       unit.Millicents `json:"price_millicents" db:"price_millicents"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`

	// Associations
	Contract            ReContract            `belongs_to:"re_contract" fk_id:"contract_id"`
	DomesticServiceArea ReDomesticServiceArea `belongs_to:"re_domestic_service_area" fk_id:"domestic_service_area_id"`
}

// ReDomesticLinehaulPrices is not required by pop and may be deleted
type ReDomesticLinehaulPrices []ReDomesticLinehaulPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReDomesticLinehaulPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.IntIsPresent{Field: r.WeightLower.Int(), Name: "WeightLower"},
		&validators.IntIsGreaterThan{Field: r.WeightLower.Int(), Name: "WeightLower", Compared: 499},
		&validators.IntIsPresent{Field: r.WeightUpper.Int(), Name: "WeightUpper"},
		&validators.IntIsGreaterThan{Field: r.WeightUpper.Int(), Name: "WeightUpper", Compared: r.WeightLower.Int()},
		&validators.IntIsGreaterThan{Field: r.MilesLower, Name: "MilesLower", Compared: -1},
		&validators.IntIsPresent{Field: r.MilesUpper, Name: "MilesUpper"},
		&validators.IntIsGreaterThan{Field: r.MilesUpper, Name: "MilesUpper", Compared: r.MilesLower},
		&validators.UUIDIsPresent{Field: r.DomesticServiceAreaID, Name: "DomesticServiceAreaID"},
		&validators.IntIsPresent{Field: r.PriceMillicents.Int(), Name: "PriceMillicents"},
		&validators.IntIsGreaterThan{Field: r.PriceMillicents.Int(), Name: "PriceMillicents", Compared: 0},
	), nil
}
