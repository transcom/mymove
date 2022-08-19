package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationAdditionalDaysSITPricer struct {
}

// NewDomesticDestinationAdditionalDaysSITPricer creates a new pricer for domestic destination additional days SIT
func NewDomesticDestinationAdditionalDaysSITPricer() services.DomesticDestinationAdditionalDaysSITPricer {
	return &domesticDestinationAdditionalDaysSITPricer{}
}

// Price determines the price for domestic destination additional days SIT
func (p domesticDestinationAdditionalDaysSITPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int, disableWeightMinimum bool) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticAdditionalDaysSIT(appCtx, models.ReServiceCodeDDASIT, contractCode, referenceDate, weight, serviceArea, numberOfDaysInSIT, disableWeightMinimum)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationAdditionalDaysSITPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	serviceAreaDest, err := getParamString(params, models.ServiceItemParamNameServiceAreaDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceAreaDest, numberOfDaysInSIT, false)
}
