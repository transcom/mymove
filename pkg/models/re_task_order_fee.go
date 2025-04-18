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

// ReTaskOrderFee model struct
type ReTaskOrderFee struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ContractYearID uuid.UUID  `json:"contract_year_id" db:"contract_year_id"`
	ServiceID      uuid.UUID  `json:"service_id" db:"service_id"`
	PriceCents     unit.Cents `json:"price_cents" db:"price_cents"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`

	//Associations
	ContractYear ReContractYear `belongs_to:"re_contract_year" fk_id:"contract_year_id"`
	Service      ReService      `belongs_to:"re_service" fk_id:"service_id"`
}

// TableName overrides the table name used by Pop.
func (r ReTaskOrderFee) TableName() string {
	return "re_task_order_fees"
}

type ReTaskOrderFees []ReTaskOrderFee

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReTaskOrderFee) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: r.ContractYearID, Name: "ContractYearID"},
		&validators.UUIDIsPresent{Field: r.ServiceID, Name: "ServiceID"},
		&validators.IntIsPresent{Field: r.PriceCents.Int(), Name: "PriceCents"},
		&validators.IntIsGreaterThan{Field: r.PriceCents.Int(), Name: "PriceCents", Compared: 0},
	), nil
}

func FetchTaskOrderFee(appCtx appcontext.AppContext, contractCode string, serviceCode ReServiceCode, date time.Time) (ReTaskOrderFee, error) {
	var taskOrderFee ReTaskOrderFee
	err := appCtx.DB().Q().
		Join("re_contract_years cy", "re_task_order_fees.contract_year_id = cy.id").
		Join("re_contracts c", "cy.contract_id = c.id").
		Join("re_services s", "re_task_order_fees.service_id = s.id").
		Where("c.code = $1", contractCode).
		Where("s.code = $2", serviceCode).
		Where("$3 between cy.start_date and cy.end_date", date).
		First(&taskOrderFee)

	if err != nil {
		return ReTaskOrderFee{}, err
	}

	return taskOrderFee, nil
}
