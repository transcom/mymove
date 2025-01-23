package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type intlOriginAdditionalDaySITPricer struct {
}

func NewIntlOriginAdditionalDaySITPricer() services.IntlOriginAdditionalDaySITPricer {
	return &intlOriginAdditionalDaySITPricer{}
}

func (p intlOriginAdditionalDaySITPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, numberOfDaysInSIT int, weight unit.Pound, perUnitCents int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceIntlAdditionalDaySIT(appCtx, models.ReServiceCodeIOASIT, contractCode, referenceDate, numberOfDaysInSIT, weight, perUnitCents)
}

func (p intlOriginAdditionalDaySITPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	numberOfDaysInSIT, err := getParamInt(params, models.ServiceItemParamNameNumberDaysSIT)
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

	return p.Price(appCtx, contractCode, referenceDate, numberOfDaysInSIT, unit.Pound(weightBilled), perUnitCents)
}
