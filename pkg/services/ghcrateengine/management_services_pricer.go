package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type managementServicesPricer struct {
	db *pop.Connection
}

// NewManagementServicesPricer creates a new pricer for management services
func NewManagementServicesPricer(db *pop.Connection) services.ManagementServicesPricer {
	return &managementServicesPricer{
		db: db,
	}
}

// Price determines the price for a management service
func (p managementServicesPricer) Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, error) {
	var taskOrderFee models.ReTaskOrderFee
	err := p.db.Q().
		Join("re_contract_years cy", "re_task_order_fees.contract_year_id = cy.id").
		Join("re_contracts c", "cy.contract_id = c.id").
		Join("re_services s", "re_task_order_fees.service_id = s.id").
		Where("c.code = $1", contractCode).
		Where("s.code = $2", models.ReServiceCodeMS).
		Where("$3 between cy.start_date and cy.end_date", mtoAvailableToPrimeAt).
		First(&taskOrderFee)

	if err != nil {
		return 0, err
	}

	return taskOrderFee.PriceCents, nil
}

// PriceUsingParams determines the price for a management service given PaymentServiceItemParams
func (p managementServicesPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), err
	}

	mtoAvailableToPrimeAt, err := getParamTime(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
	if err != nil {
		return unit.Cents(0), err
	}

	return p.Price(contractCode, mtoAvailableToPrimeAt)
}
