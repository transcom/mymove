package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type counselingServicesPricer struct {
}

// NewCounselingServicesPricer creates a new pricer for counseling services
func NewCounselingServicesPricer() services.CounselingServicesPricer {
	return &counselingServicesPricer{}
}

// Price determines the price for a counseling service
func (p counselingServicesPricer) Price(appCtx appcontext.AppContext, contractCode string, mtoAvailableToPrimeAt time.Time) (unit.Cents, services.PricingDisplayParams, error) {
	taskOrderFee, err := fetchTaskOrderFee(appCtx, contractCode, models.ReServiceCodeCS, mtoAvailableToPrimeAt)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch task order fee: %w", err)
	}

	displayPriceParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(taskOrderFee.PriceCents),
		},
	}
	return taskOrderFee.PriceCents, displayPriceParams, nil
}

// PriceUsingParams determines the price for a counseling service given PaymentServiceItemParams
func (p counselingServicesPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	mtoAvailableToPrimeAt, err := getParamTime(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, mtoAvailableToPrimeAt)
}
