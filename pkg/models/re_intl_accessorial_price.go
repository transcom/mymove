package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// Market represents the market for an international move
type Market string

func (m Market) String() string {
	return string(m)
}

// This lists available markets for international accessorial pricing
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

// TableName overrides the table name used by Pop.
func (r ReIntlAccessorialPrice) TableName() string {
	return "re_intl_accessorial_prices"
}

type ReIntlAccessorialPrices []ReIntlAccessorialPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReIntlAccessorialPrice) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.StringIsPresent{Field: r.Market.String(), Name: "Market"},
		&validators.StringInclusion{Field: r.Market.String(), Name: "Market", List: validMarkets},
		&validators.IntIsGreaterThan{Field: r.PerUnitCents.Int(), Name: "PerUnitCents", Compared: -1},
	), nil
}

func (m Market) FullString() string {
	switch m {
	case MarketConus:
		return "CONUS"
	case MarketOconus:
		return "OCONUS"
	default:
		return ""
	}
}
