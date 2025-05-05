package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
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

// TableName overrides the table name used by Pop.
func (r ReDomesticAccessorialPrice) TableName() string {
	return "re_domestic_accessorial_prices"
}

type ReDomesticAccessorialPrices []ReDomesticAccessorialPrice

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReDomesticAccessorialPrice) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractID, Name: "ContractID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.IntIsGreaterThan{Field: r.ServicesSchedule, Name: "ServicesSchedule", Compared: 0},
		&validators.IntIsLessThan{Field: r.ServicesSchedule, Name: "ServicesSchedule", Compared: 4},
		&validators.IntIsPresent{Field: r.PerUnitCents.Int(), Name: "PerUnitCents"},
		&validators.IntIsGreaterThan{Field: r.PerUnitCents.Int(), Name: "PerUnitCents", Compared: 0},
	), nil
}

func FetchAccessorialPrice(appCtx appcontext.AppContext, contractCode string, serviceCode ReServiceCode, schedule int) (ReDomesticAccessorialPrice, error) {
	var domAccessorialPrice ReDomesticAccessorialPrice
	err := appCtx.DB().Q().
		Join("re_services", "service_id = re_services.id").
		Join("re_contracts", "re_contracts.id = re_domestic_accessorial_prices.contract_id").
		Where("re_contracts.code = $1", contractCode).
		Where("re_services.code = $2", serviceCode).
		Where("services_schedule = $3", schedule).
		First(&domAccessorialPrice)

	if err != nil {
		return ReDomesticAccessorialPrice{}, err
	}

	return domAccessorialPrice, nil
}
