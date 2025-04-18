package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// ReIntlOtherPrice is the ghc rate engine international price
type ReIntlOtherPrice struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ContractID    uuid.UUID  `json:"contract_id" db:"contract_id"`
	ServiceID     uuid.UUID  `json:"service_id" db:"service_id"`
	RateAreaID    uuid.UUID  `json:"rate_area_id" db:"rate_area_id"`
	IsPeakPeriod  bool       `json:"is_peak_period" db:"is_peak_period"`
	PerUnitCents  unit.Cents `json:"per_unit_cents" db:"per_unit_cents"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	IsLess50Miles *bool      `json:"is_less_50_miles" db:"is_less_50_miles"`

	// Associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	Service  ReService  `belongs_to:"re_service" fk_id:"service_id"`
	RateArea ReRateArea `belongs_to:"re_rate_area" fk_id:"rate_area_id"`
}

// TableName overrides the table name used by Pop.
func (r ReIntlOtherPrice) TableName() string {
	return "re_intl_other_prices"
}

// ReIntlOtherPrices is a slice of ReIntlOtherPrice
type ReIntlOtherPrices []ReIntlOtherPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReIntlOtherPrice) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.UUIDIsPresent{Field: r.RateAreaID, Name: "RateAreaID"},
		&validators.IntIsGreaterThan{Field: r.PerUnitCents.Int(), Name: "PerUnitCents", Compared: -1},
	), nil
}

// fetches a row from re_intl_other_prices using passed in parameters
// gets the rate_area_id & is_peak_period based on values provided
func FetchReIntlOtherPrice(db *pop.Connection, addressID uuid.UUID, serviceID uuid.UUID, contractID uuid.UUID, referenceDate *time.Time) (*ReIntlOtherPrice, error) {
	if addressID != uuid.Nil && serviceID != uuid.Nil && contractID != uuid.Nil && referenceDate != nil {
		// need to get the rate area first
		rateAreaID, err := FetchRateAreaID(db, addressID, &serviceID, contractID)
		if err != nil {
			return nil, fmt.Errorf("error fetching rate area id for shipment ID: %s, service ID %s, and contract ID: %s: %s", addressID, serviceID, contractID, err)
		}

		var isPeakPeriod bool
		err = db.RawQuery("SELECT is_peak_period($1)", referenceDate).First(&isPeakPeriod)
		if err != nil {
			return nil, fmt.Errorf("error checking if date is peak period with date: %s: %s", contractID, err)
		}

		var reIntlOtherPrice ReIntlOtherPrice
		err = db.Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", rateAreaID).
			First(&reIntlOtherPrice)
		if err != nil {
			return nil, fmt.Errorf("error fetching row from re_int_other_prices using rateAreaID %s, service ID %s, and contract ID: %s: %s", rateAreaID, serviceID, contractID, err)
		}

		return &reIntlOtherPrice, nil
	}

	return nil, fmt.Errorf("error value from re_intl_other_prices - required parameters not provided")
}
