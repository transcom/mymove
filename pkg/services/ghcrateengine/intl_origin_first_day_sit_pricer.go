package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type intlOriginFirstDaySITPricer struct {
}

func NewIntlOriginFirstDaySITPricer() services.IntlOriginFirstDaySITPricer {
	return &intlOriginFirstDaySITPricer{}
}

func (p intlOriginFirstDaySITPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, perUnitCents int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceIntlFirstDaySIT(appCtx, models.ReServiceCodeIOFSIT, contractCode, referenceDate, weight, perUnitCents)
}

func (p intlOriginFirstDaySITPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	perUnitCents, err := getParamInt(params, models.ServiceItemParamNamePerUnitCents)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), perUnitCents)
}
