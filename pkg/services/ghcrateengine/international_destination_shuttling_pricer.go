package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type internationalDestinationShuttlingPricer struct {
}

// NewInternationalDestinationShuttlingPricer creates a new pricer for international destination first day SIT
func NewInternationalDestinationShuttlingPricer() services.InternationalDestinationShuttlingPricer {
	return &internationalDestinationShuttlingPricer{}
}

// Price determines the price for international destination first day SIT
func (p internationalDestinationShuttlingPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, market models.Market) (unit.Cents, services.PricingDisplayParams, error) {
	return priceInternationalShuttling(appCtx, models.ReServiceCodeIDSHUT, contractCode, referenceDate, weight, market)
}

// PriceUsingParams determines the price for international destination first day SIT given PaymentServiceItemParams
func (p internationalDestinationShuttlingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	market, err := getParamMarket(params, models.ServiceItemParamNameMarketDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), market)
}
