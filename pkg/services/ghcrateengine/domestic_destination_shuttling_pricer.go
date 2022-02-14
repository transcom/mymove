package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationShuttlingPricer struct {
}

// NewDomesticDestinationShuttlingPricer creates a new pricer for domestic destination first day SIT
func NewDomesticDestinationShuttlingPricer() services.DomesticDestinationShuttlingPricer {
	return &domesticDestinationShuttlingPricer{}
}

// Price determines the price for domestic destination first day SIT
func (p domesticDestinationShuttlingPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticShuttling(appCtx, models.ReServiceCodeDDSHUT, contractCode, referenceDate, weight, serviceSchedule)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationShuttlingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceScheduleDestination, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceScheduleDestination)
}
