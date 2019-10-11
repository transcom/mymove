package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ReIntlAccessorialPrice model struct
type ReIntlAccessorialPrice struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ContractID   uuid.UUID `json:"contract_id" db:"contract_id"`
	ServiceID    uuid.UUID `json:"service_id" db:"service_id"`
	Market       string    `json:"market" db:"market"`
	PerUnitCents int       `json:"per_unit_cents" db:"per_unit_cents"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	//associations
	Contract ReContract `belongs_to:"re_contract"`
	Service  ReService  `belongs_to:"re_service"`
}

// ReIntlAccessorialPrices is not required by pop and may be deleted
type ReIntlAccessorialPrices []ReIntlAccessorialPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReIntlAccessorialPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	validMarkets := []string{"C", "O"}
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceTypeID"},
		&validators.StringIsPresent{Field: r.Market, Name: "Market"},
		&validators.StringInclusion{Field: r.Market, Name: "Market", List: validMarkets},
		&validators.IntIsPresent{Field: r.PerUnitCents, Name: "PerUnitCents"},
	), nil
}
