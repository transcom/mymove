package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// ReIntlOtherPrice is the ghc rate engine international price
type ReIntlOtherPrice struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ContractID   uuid.UUID  `json:"contract_id" db:"contract_id"`
	ServiceID    uuid.UUID  `json:"service_id" db:"service_id"`
	RateAreaID   uuid.UUID  `json:"rate_area_id" db:"rate_area_id"`
	IsPeakPeriod bool       `json:"is_peak_period" db:"is_peak_period"`
	PerUnitCents unit.Cents `json:"per_unit_cents" db:"per_unit_cents"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	// Associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	Service  ReService  `belongs_to:"re_service" fk_id:"service_id"`
	RateArea ReRateArea `belongs_to:"re_rate_area" fk_id:"rate_area_id"`
}

// ReIntlOtherPrices is a slice of ReIntlOtherPrice
type ReIntlOtherPrices []ReIntlOtherPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *ReIntlOtherPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: p.ServiceID, Name: "ServiceID"},
		&validators.UUIDIsPresent{Field: p.RateAreaID, Name: "RateAreaID"},
		&validators.IntIsGreaterThan{Field: p.PerUnitCents.Int(), Name: "PerUnitCents", Compared: 0},
	), nil
}
