package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// ReDomesticAccessorialPrice model struct
type ReDomesticAccessorialPrice struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ContractID       uuid.UUID  `json:"contract_id" db:"contract_id"`
	ServiceID        uuid.UUID  `json:"service_id" db:"service_id"`
	ServicesSchedule int        `json:"services_schedule" db:"services_schedule"`
	PerUnitCents     unit.Cents `json:"per_unit_cents" db:"per_unit_cents"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`

	//associations
	Contract ReContract `belongs_to:"re_contract" fk_id:"contract_id"`
	Service  ReService  `belongs_to:"re_service" fk_id:"service_id"`
}

// ReDomesticAccessorialPrices is not required by pop and may be deleted
type ReDomesticAccessorialPrices []ReDomesticAccessorialPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReDomesticAccessorialPrice) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.IntIsGreaterThan{Field: r.ServicesSchedule, Name: "ServicesSchedule", Compared: 0},
		&validators.IntIsLessThan{Field: r.ServicesSchedule, Name: "ServicesSchedule", Compared: 4},
		&validators.IntIsPresent{Field: r.PerUnitCents.Int(), Name: "PerUnitCents"},
		&validators.IntIsGreaterThan{Field: r.PerUnitCents.Int(), Name: "PerUnitCents", Compared: 0},
	), nil
}
