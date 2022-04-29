package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

//Market represents the market for an international move
type Market string

func (m Market) String() string {
	return string(m)
}

//This lists available markets for international accessorial pricing
const (
	MarketConus  Market = "C"
	MarketOconus Market = "O"
)

var validMarkets = []string{
	string(MarketConus),
	string(MarketOconus),
}

// ReIntlAccessorialPrice model struct
type ReIntlAccessorialPrice struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ContractID   uuid.UUID  `json:"contract_id" db:"contract_id"`
	ServiceID    uuid.UUID  `json:"service_id" db:"service_id"`
	Market       Market     `json:"market" db:"market"`
	PerUnitCents unit.Cents `json:"per_unit_cents" db:"per_unit_cents"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	//associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	Service  ReService  `belongs_to:"re_service" fk_id:"service_id"`
}

// ReIntlAccessorialPrices is not required by pop and may be deleted
type ReIntlAccessorialPrices []ReIntlAccessorialPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReIntlAccessorialPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.StringIsPresent{Field: r.Market.String(), Name: "Market"},
		&validators.StringInclusion{Field: r.Market.String(), Name: "Market", List: validMarkets},
		&validators.IntIsPresent{Field: r.PerUnitCents.Int(), Name: "PerUnitCents"},
		&validators.IntIsGreaterThan{Field: r.PerUnitCents.Int(), Name: "PerUnitCents", Compared: 0},
	), nil
}
