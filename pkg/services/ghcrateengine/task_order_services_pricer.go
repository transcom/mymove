package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type taskOrderServicesPricer struct {
	db          *pop.Connection
	serviceCode models.ReServiceCode
}

// NewTaskOrderServicesPricer creates a new pricer for task order services given explicit parameter values
func NewTaskOrderServicesPricer(db *pop.Connection, serviceCode models.ReServiceCode) services.TaskOrderServicesPricer {
	return &taskOrderServicesPricer{
		db:          db,
		serviceCode: serviceCode,
	}
}

// Price determines the price for a task order service
func (p taskOrderServicesPricer) Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, error) {
	if p.serviceCode != models.ReServiceCodeMS && p.serviceCode != models.ReServiceCodeCS {
		return 0, fmt.Errorf("invalid service code: %s", p.serviceCode)
	}

	var taskOrderFee models.ReTaskOrderFee
	err := p.db.Q().
		Join("re_contract_years cy", "re_task_order_fees.contract_year_id = cy.id").
		Join("re_contracts c", "cy.contract_id = c.id").
		Join("re_services s", "re_task_order_fees.service_id = s.id").
		Where("c.code = $1", contractCode).
		Where("s.code = $2", p.serviceCode).
		Where("$3 between cy.start_date and cy.end_date", mtoAvailableToPrimeAt).
		First(&taskOrderFee)

	if err != nil {
		return 0, err
	}

	return taskOrderFee.PriceCents, nil
}

// PriceUsingParams determines the price for a task order service given PaymentServiceItemParams
func (p taskOrderServicesPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
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
