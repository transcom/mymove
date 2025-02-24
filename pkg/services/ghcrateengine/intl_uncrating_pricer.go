package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type intlUncratingPricer struct {
}

// NewIntlUncratingPricer creates a new pricer for international uncrating
func NewIntlUncratingPricer() services.IntlUncratingPricer {
	return &intlUncratingPricer{}
}

// Price determines the price for international uncrating
func (p intlUncratingPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, billedCubicFeet unit.CubicFeet, market models.Market) (unit.Cents, services.PricingDisplayParams, error) {
	return priceIntlCratingUncrating(appCtx, models.ReServiceCodeIUCRT, contractCode, referenceDate, billedCubicFeet, false, 0, false, market)
}

// PriceUsingParams determines the price for international uncrating given PaymentServiceItemParams
func (p intlUncratingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetFloat, err := getParamFloat(params, models.ServiceItemParamNameCubicFeetBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetBilled := unit.CubicFeet(cubicFeetFloat)

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	market, err := getParamMarket(params, models.ServiceItemParamNameMarketDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, cubicFeetBilled, market)
}
