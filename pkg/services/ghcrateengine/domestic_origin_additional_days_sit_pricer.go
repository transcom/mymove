package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginAdditionalDaysSITPricer struct {
}

// NewDomesticOriginAdditionalDaysSITPricer creates a new pricer for domestic origin additional days SIT
func NewDomesticOriginAdditionalDaysSITPricer() services.DomesticOriginAdditionalDaysSITPricer {
	return &domesticOriginAdditionalDaysSITPricer{}
}

// Price determines the price for domestic origin additional days SIT
func (p domesticOriginAdditionalDaysSITPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticAdditionalDaysSIT(appCtx, models.ReServiceCodeDOASIT, contractCode, referenceDate, weight, serviceArea, numberOfDaysInSIT)
}

// PriceUsingParams determines the price for domestic origin first day SIT given PaymentServiceItemParams
func (p domesticOriginAdditionalDaysSITPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	numberOfDaysInSIT, err := getParamInt(params, models.ServiceItemParamNameNumberDaysSIT)
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

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceAreaOrigin, numberOfDaysInSIT)
}
