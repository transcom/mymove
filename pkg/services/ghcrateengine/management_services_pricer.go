package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"

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
	taskOrderFee, err := fetchTaskOrderFee(p.db, contractCode, models.ReServiceCodeMS, mtoAvailableToPrimeAt)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch task order fee: %w", err)
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
