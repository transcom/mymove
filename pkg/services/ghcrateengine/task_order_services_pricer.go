package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type taskOrderServicesPricer struct {
	db                    *pop.Connection
	contractCode          string
	serviceCode           models.ReServiceCode
	mtoAvailableToPrimeAt time.Time
}

// NewTaskOrderServicesPricer creates a new pricer for task order services given explicit parameter values
func NewTaskOrderServicesPricer(db *pop.Connection, contractCode string, serviceCode models.ReServiceCode, mtoAvailableToPrimeAt time.Time) Pricer {
	return &taskOrderServicesPricer{
		db:                    db,
		contractCode:          contractCode,
		serviceCode:           serviceCode,
		mtoAvailableToPrimeAt: mtoAvailableToPrimeAt,
	}
}

// NewTaskOrderServicesPricerFromParams creates a new pricer for task order services based upon PaymentServiceItemParams
func NewTaskOrderServicesPricerFromParams(db *pop.Connection, serviceCode models.ReServiceCode, params models.PaymentServiceItemParams) (Pricer, error) {
	pricer := &taskOrderServicesPricer{
		db:           db,
		contractCode: testdatagen.DefaultContractCode, // TODO: What to do about contract code?
		serviceCode:  serviceCode,
	}

	var err error
	pricer.mtoAvailableToPrimeAt, err = getParamTime(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
	if err != nil {
		return nil, err
	}

	return pricer, nil
}

// Price determines the price for a task order service
func (p taskOrderServicesPricer) Price() (unit.Cents, error) {
	if p.serviceCode != models.ReServiceCodeMS && p.serviceCode != models.ReServiceCodeCS {
		return 0, fmt.Errorf("invalid service code: %s", p.serviceCode)
	}

	var taskOrderFee models.ReTaskOrderFee
	err := p.db.Q().
		Join("re_contract_years cy", "re_task_order_fees.contract_year_id = cy.id").
		Join("re_contracts c", "cy.contract_id = c.id").
		Join("re_services s", "re_task_order_fees.service_id = s.id").
		Where("c.code = $1", p.contractCode).
		Where("s.code = $2", p.serviceCode).
		Where("$3 between cy.start_date and cy.end_date", p.mtoAvailableToPrimeAt).
		First(&taskOrderFee)

	if err != nil {
		return 0, err
	}

	return taskOrderFee.PriceCents, nil
}
