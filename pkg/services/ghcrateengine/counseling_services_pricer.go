package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type counselingServicesPricer struct {
	db *pop.Connection
}

// NewCounselingServicesPricer creates a new pricer for counseling services
func NewCounselingServicesPricer(db *pop.Connection) services.CounselingServicesPricer {
	return &counselingServicesPricer{
		db: db,
	}
}

// Price determines the price for a counseling service
func (p counselingServicesPricer) Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, error) {
	var taskOrderFee models.ReTaskOrderFee
	err := p.db.Q().
		Join("re_contract_years cy", "re_task_order_fees.contract_year_id = cy.id").
		Join("re_contracts c", "cy.contract_id = c.id").
		Join("re_services s", "re_task_order_fees.service_id = s.id").
		Where("c.code = $1", contractCode).
		Where("s.code = $2", models.ReServiceCodeCS).
		Where("$3 between cy.start_date and cy.end_date", mtoAvailableToPrimeAt).
		First(&taskOrderFee)

	if err != nil {
		return 0, err
	}

	return taskOrderFee.PriceCents, nil
}

// PriceUsingParams determines the price for a counseling service given PaymentServiceItemParams
func (p counselingServicesPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
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
