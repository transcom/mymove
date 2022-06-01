package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ReIntlPrice is the ghc rate engine international price
type ReIntlPrice struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	ContractID            uuid.UUID  `json:"contract_id" db:"contract_id"`
	ServiceID             uuid.UUID  `json:"service_id" db:"service_id"`
	OriginRateAreaID      uuid.UUID  `json:"origin_rate_area_id" db:"origin_rate_area_id"`
	DestinationRateAreaID uuid.UUID  `json:"destination_rate_area_id" db:"destination_rate_area_id"`
	IsPeakPeriod          bool       `json:"is_peak_period" db:"is_peak_period"`
	PerUnitCents          unit.Cents `json:"per_unit_cents" db:"per_unit_cents"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`

	// Associations
	Contract            ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	Service             ReService  `belongs_to:"re_service" fk_id:"service_id"`
	OriginRateArea      ReRateArea `belongs_to:"re_rate_area" fk_id:"origin_rate_area_id"`
	DestinationRateArea ReRateArea `belongs_to:"re_rate_area" fk_id:"destination_rate_area_id"`
}

// ReIntlPrices is a slice of ReIntlPrice objects
type ReIntlPrices []ReIntlPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (p *ReIntlPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: p.ServiceID, Name: "ServiceID"},
		&validators.UUIDIsPresent{Field: p.OriginRateAreaID, Name: "OriginRateAreaID"},
		&validators.UUIDIsPresent{Field: p.DestinationRateAreaID, Name: "DestinationRateAreaID"},
		&validators.IntIsGreaterThan{Field: p.PerUnitCents.Int(), Name: "PerUnitCents", Compared: 0},
	), nil
}
