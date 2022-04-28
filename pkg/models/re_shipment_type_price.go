package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ReShipmentTypePrice model struct
type ReShipmentTypePrice struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ContractID uuid.UUID `json:"contract_id" db:"contract_id"`
	ServiceID  uuid.UUID `json:"service_id" db:"service_id"`
	Market     Market    `json:"market" db:"market"`
	Factor     float64   `json:"factor" db:"factor"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	//Associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	Service  ReService  `belongs_to:"re_service" fk_id:"service_id"`
}

// ReShipmentTypePrices is not required by pop and may be deleted
type ReShipmentTypePrices []ReShipmentTypePrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReShipmentTypePrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.StringIsPresent{Field: r.Market.String(), Name: "Market"},
		&validators.StringInclusion{Field: r.Market.String(), Name: "Market", List: validMarkets},
		&Float64IsPresent{Field: r.Factor, Name: "Factor"},
		&Float64IsGreaterThan{Field: r.Factor, Name: "Factor", Compared: 0},
	), nil
}
