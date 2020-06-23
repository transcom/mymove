package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

func fetchTaskOrderFee(db *pop.Connection, contractCode string, serviceCode models.ReServiceCode, mtoAvailableToPrimeAt time.Time) (models.ReTaskOrderFee, error) {
	var taskOrderFee models.ReTaskOrderFee
	err := db.Q().
		Join("re_contract_years cy", "re_task_order_fees.contract_year_id = cy.id").
		Join("re_contracts c", "cy.contract_id = c.id").
		Join("re_services s", "re_task_order_fees.service_id = s.id").
		Where("c.code = $1", contractCode).
		Where("s.code = $2", serviceCode).
		Where("$3 between cy.start_date and cy.end_date", mtoAvailableToPrimeAt).
		First(&taskOrderFee)

	if err != nil {
		return models.ReTaskOrderFee{}, err
	}

	return taskOrderFee, nil
}
