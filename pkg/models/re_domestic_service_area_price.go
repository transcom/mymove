package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// ReDomesticServiceAreaPrice represents a domestic service area price based on date, service area, etc.
type ReDomesticServiceAreaPrice struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	ContractID            uuid.UUID  `json:"contract_id" db:"contract_id"`
	ServiceID             uuid.UUID  `json:"service_id" db:"service_id"`
	IsPeakPeriod          bool       `json:"is_peak_period" db:"is_peak_period"`
	DomesticServiceAreaID uuid.UUID  `json:"domestic_service_area_id" db:"domestic_service_area_id"`
	PriceCents            unit.Cents `json:"price_cents" db:"price_cents"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`

	// Associations
	Contract            ReContract            `belongs_to:"re_contract" fk_id:"contract_id"`
	Service             ReService             `belongs_to:"re_service" fk_id:"service_id"`
	DomesticServiceArea ReDomesticServiceArea `belongs_to:"re_domestic_service_area" fk_id:"domestic_service_area_id"`
}

// ReDomesticServiceAreaPrices is not required by pop and may be deleted
type ReDomesticServiceAreaPrices []ReDomesticServiceAreaPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReDomesticServiceAreaPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.UUIDIsPresent{Field: r.DomesticServiceAreaID, Name: "DomesticServiceAreaID"},
		&validators.IntIsPresent{Field: r.PriceCents.Int(), Name: "PriceCents"},
		&validators.IntIsGreaterThan{Field: r.PriceCents.Int(), Name: "PriceCents", Compared: 0},
	), nil
}
