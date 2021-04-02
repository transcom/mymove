package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"

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
func (p counselingServicesPricer) Price(contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, services.PricingDisplayParams, error) {
	taskOrderFee, err := fetchTaskOrderFee(p.db, contractCode, models.ReServiceCodeCS, mtoAvailableToPrimeAt)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch task order fee: %w", err)
	}

	return taskOrderFee.PriceCents, nil, nil
}

// PriceUsingParams determines the price for a counseling service given PaymentServiceItemParams
func (p counselingServicesPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	mtoAvailableToPrimeAt, err := getParamTime(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(contractCode, mtoAvailableToPrimeAt)
}
