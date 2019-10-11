package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ReTaskOrderPrice model struct
type ReTaskOrderFee struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ContractYearID uuid.UUID `json:"contract_year_id" db:"contract_year_id"`
	ServiceID      uuid.UUID `json:"service_type_id" db:"service_type_id"`
	PriceCents     int       `json:"price_cents" db:"price_cents"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`

	//Associations
	ContractYear ReContractYear `belongs_to:"re_contract_year"`
	Service      ReService      `belongs_to:"re_service"`
}

// ReTaskOrderPrices is not required by pop and may be deleted
type ReTaskOrderFees []ReTaskOrderFee

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReTaskOrderFee) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractYearID, Name: "ContractYearID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceTypeID"},
		&validators.IntIsPresent{Field: r.PriceCents, Name: "PriceCents"},
	), nil
}
