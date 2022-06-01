package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// ReDomesticOtherPrice represents a domestic service area price based on date, service area, etc.
type ReDomesticOtherPrice struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ContractID   uuid.UUID  `json:"contract_id" db:"contract_id"`
	ServiceID    uuid.UUID  `json:"service_id" db:"service_id"`
	IsPeakPeriod bool       `json:"is_peak_period" db:"is_peak_period"`
	Schedule     int        `json:"schedule" db:"schedule"`
	PriceCents   unit.Cents `json:"price_cents" db:"price_cents"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	// Associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	Service  ReService  `belongs_to:"re_service" fk_id:"service_id"`
}

// ReDomesticOtherPrices is not required by pop and may be deleted
type ReDomesticOtherPrices []ReDomesticOtherPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReDomesticOtherPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.IntIsGreaterThan{Field: r.Schedule, Name: "Schedule", Compared: 0},
		&validators.IntIsLessThan{Field: r.Schedule, Name: "Schedule", Compared: 4},
		&validators.IntIsPresent{Field: r.PriceCents.Int(), Name: "PriceCents"},
		&validators.IntIsGreaterThan{Field: r.PriceCents.Int(), Name: "PriceCents", Compared: 0},
	), nil
}
