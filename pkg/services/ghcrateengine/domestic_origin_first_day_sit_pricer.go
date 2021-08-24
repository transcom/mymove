package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginFirstDaySITPricer struct {
}

// NewDomesticOriginFirstDaySITPricer creates a new pricer for domestic origin first day SIT
func NewDomesticOriginFirstDaySITPricer() services.DomesticOriginFirstDaySITPricer {
	return &domesticOriginFirstDaySITPricer{}
}

// Price determines the price for domestic origin first day SIT
func (p domesticOriginFirstDaySITPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticFirstDaySIT(appCtx, models.ReServiceCodeDOFSIT, contractCode, requestedPickupDate, weight, serviceArea)
}

// PriceUsingParams determines the price for domestic origin first day SIT given PaymentServiceItemParams
func (p domesticOriginFirstDaySITPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceAreaOrigin)
}
