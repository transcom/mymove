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
func (p domesticOriginFirstDaySITPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, disableWeightMinimum bool) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticFirstDaySIT(appCtx, models.ReServiceCodeDOFSIT, contractCode, referenceDate, weight, serviceArea, disableWeightMinimum)
}

// PriceUsingParams determines the price for domestic origin first day SIT given PaymentServiceItemParams
func (p domesticOriginFirstDaySITPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceAreaOrigin, false)
}
