package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationFirstDaySITPricer struct {
}

// NewDomesticDestinationFirstDaySITPricer creates a new pricer for domestic destination first day SIT
func NewDomesticDestinationFirstDaySITPricer() services.DomesticDestinationFirstDaySITPricer {
	return &domesticDestinationFirstDaySITPricer{}
}

// Price determines the price for domestic destination first day SIT
func (p domesticDestinationFirstDaySITPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, disableWeightMinimum bool) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticFirstDaySIT(appCtx, models.ReServiceCodeDDFSIT, contractCode, referenceDate, weight, serviceArea, disableWeightMinimum)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationFirstDaySITPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
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

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceAreaDest, true)
}
