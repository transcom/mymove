package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type internationalOriginShuttlingPricer struct {
}

// NewInternationalOriginShuttlingPricer creates a new pricer for international origin first day SIT
func NewInternationalOriginShuttlingPricer() services.InternationalOriginShuttlingPricer {
	return &internationalOriginShuttlingPricer{}
}

// Price determines the price for international origin first day SIT
func (p internationalOriginShuttlingPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceInternationalShuttling(appCtx, models.ReServiceCodeDOSHUT, contractCode, referenceDate, weight, serviceSchedule)
}

// PriceUsingParams determines the price for international origin first day SIT given PaymentServiceItemParams
func (p internationalOriginShuttlingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceScheduleOrigin, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceScheduleOrigin)
}
